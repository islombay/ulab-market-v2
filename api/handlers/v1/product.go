package handlersv1

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage"
	"app/storage/filestore"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"net/http"
	"strings"
	"time"
)

// CreateProduct
// @id createProduct
// @router /api/product [post]
// @summary create product
// @description create product
// @security ApiKeyAuth
// @accept multipart/form-data
// @tags product
// @produce json
// @param createProduct formData models_v1.CreateProduct true "create product"
// @param main_image formData file false "Main image"
// @param image_files formData []file false "Image files (multiple)"
// @param video_files formData []file false "Video files (multiple)"
// @Success 200 {object} models.Product "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid / status invalid"
// @Failure 404 {object} models_v1.Response "Category not found / Brand not found"
// @failure 409 {object} models_v1.Response "Articul already found"
// @Failure 413 {object} models_v1.Response "Image size is big / Video size is big / Image count too many / Video count too many"
// @Failure 415 {object} models_v1.Response "Image type is not supported / Video type is not supported"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) CreateProduct(c *gin.Context) {
	var m models_v1.CreateProduct
	if err := c.Bind(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		v1.log.Error("bad request", logs.Error(err))
		return
	}

	_, err := v1.storage.Product().GetByArticul(context.Background(), m.Articul)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get product by articul", logs.Error(err))
			return
		}
	} else {
		v1.error(c, status.StatusAlreadyExists)
		return
	}

	m.Status = strings.ToLower(m.Status)
	if !(m.Status == helper.ProductStatusActive || m.Status == helper.ProductStatusInactive) {
		v1.error(c, status.StatusProductStatusInvalid)
		return
	}

	if m.CategoryID != "" {
		if !helper.IsValidUUID(m.CategoryID) {
			v1.error(c, status.StatusBadUUID)
			return
		}
		if _, err := v1.storage.Category().GetByID(context.Background(), m.CategoryID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.error(c, status.StatusCategoryNotFound)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get category by id", logs.Error(err),
				logs.String("cid", m.CategoryID),
			)
			return
		}
	}

	if m.BrandID != "" {
		if !helper.IsValidUUID(m.BrandID) {
			v1.error(c, status.StatusBadUUID)
			return
		}
		if _, err := v1.storage.Brand().GetByID(context.Background(), m.BrandID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.error(c, status.StatusBrandNotFound)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get brand by id", logs.Error(err),
				logs.String("bid", m.BrandID),
			)
			return
		}
	}

	pr := models.Product{
		ID:            uuid.NewString(),
		Articul:       m.Articul,
		NameUz:        m.NameUz,
		NameRu:        m.NameRu,
		DescriptionUz: m.DescriptionUz,
		DescriptionRu: m.DescriptionRu,
		IncomePrice:   m.IncomePrice,
		OutcomePrice:  m.OutcomePrice,
		Quantity:      m.Quantity,
		CategoryID:    models.GetStringAddress(m.CategoryID),
		BrandID:       models.GetStringAddress(m.BrandID),
		Status:        m.Status,
		MainImage:     nil,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if m.MainImage != nil {
		if _, err, msg := helper.IsValidImage(m.MainImage); err != nil {
			if errors.Is(err, helper.ErrInvalidImageType) {
				v1.error(c, status.StatusImageTypeUnkown)
				v1.log.Error("got invalid image extension", logs.String("got", msg))
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not check file for validity")
			return
		}
		if m.MainImage.Size > v1.cfg.Media.ProductPhotoMaxSize {
			v1.error(c, status.StatusProductMainImageMaxSizeExceed)
			return
		}
		url, err := v1.filestore.Create(m.MainImage, filestore.FolderProduct, pr.ID)
		if err != nil {
			v1.error(c, status.StatusInternal)
			return
		}
		pr.MainImage = models.GetStringAddress(url)
	}

	if err := v1.storage.Product().CreateProduct(context.Background(), pr); err != nil {
		if pr.MainImage != nil {
			v1.filestore.DeleteFile(*pr.MainImage)
		}
		if errors.Is(err, storage.ErrNotAffected) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("got not affected for creating product")
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not create product", logs.Error(err))
		return
	}

	// add image adding
	if len(m.ImageFiles) > v1.cfg.Media.ProductPhotoMaxCount {
		v1.error(c, status.StatusProductPhotoMaxCount)
		return
	}
	var imageFilesStopped struct {
		Status *status.Status
	}
	imageFiles := make([]models.ProductMediaFiles, len(m.ImageFiles))
	for ind, i := range m.ImageFiles {
		if _, err, msg := helper.IsValidImage(i); err != nil {
			if errors.Is(err, helper.ErrInvalidImageType) {
				imageFilesStopped.Status = &status.StatusImageTypeUnkown
				v1.log.Error("got invalid image extension", logs.String("got", msg))
				break
			}
			imageFilesStopped.Status = &status.StatusInternal
			v1.log.Error("could not check file for validity")
			break
		}
		if i.Size > v1.cfg.Media.ProductPhotoMaxSize {
			imageFilesStopped.Status = &status.StatusImageMaxSizeExceed
			break
		}
		id := uuid.NewString()
		url, err := v1.filestore.Create(i, filestore.FolderProduct, id)
		if err != nil {
			imageFilesStopped.Status = &status.StatusInternal
			break
		}
		if err := v1.storage.Product().CreateProductImageFile(context.Background(),
			id, pr.ID, url,
		); err != nil {
			imageFilesStopped.Status = &status.StatusInternal
			v1.log.Error("could not create product image file in db", logs.Error(err),
				logs.String("pid", pr.ID),
			)
			break
		}
		imageFiles[ind] = models.ProductMediaFiles{
			ID:        id,
			ProductID: pr.ID,
			MediaFile: v1.filestore.GetURL(url),
		}
	}
	if imageFilesStopped.Status != nil {
		v1.error(c, *imageFilesStopped.Status)
		for _, img := range imageFiles {
			v1.filestore.DeleteFile(img.MediaFile)
		}
		if pr.MainImage != nil {
			v1.filestore.DeleteFile(*pr.MainImage)
		}
		v1.storage.Product().DeleteProductByID(context.Background(), pr.ID)
		return
	}

	// add video adding
	if len(m.VideoFiles) > v1.cfg.Media.ProductVideoMaxCount {
		v1.error(c, status.StatusProductVideoMaxCount)
		return
	}
	videoFiles := make([]models.ProductMediaFiles, len(m.VideoFiles))
	for ind, i := range m.VideoFiles {
		if _, err, msg := helper.IsValidVideo(i); err != nil {
			if errors.Is(err, helper.ErrInvalidVideoType) {
				imageFilesStopped.Status = &status.StatusVideoTypeUnkown
				v1.log.Error("got invalid video extension", logs.String("got", msg))
				break
			}
			imageFilesStopped.Status = &status.StatusInternal
			v1.log.Error("could not check file for validity")
			break
		}
		id := uuid.NewString()
		url, err := v1.filestore.Create(i, filestore.FolderProduct, id)
		if err != nil {
			imageFilesStopped.Status = &status.StatusInternal
			for _, img := range imageFiles {
				v1.filestore.DeleteFile(img.MediaFile)
			}
			break
		}
		if err := v1.storage.Product().CreateProductVideoFile(context.Background(),
			id, pr.ID, url,
		); err != nil {
			imageFilesStopped.Status = &status.StatusInternal
			v1.log.Error("could not create product video file in db", logs.Error(err),
				logs.String("pid", pr.ID),
			)

			for _, img := range imageFiles {
				v1.filestore.DeleteFile(img.MediaFile)
			}

			break
		}
		videoFiles[ind] = models.ProductMediaFiles{
			ID:        id,
			ProductID: pr.ID,
			MediaFile: v1.filestore.GetURL(url),
		}
	}
	if imageFilesStopped.Status != nil {
		v1.error(c, *imageFilesStopped.Status)
		for _, img := range imageFiles {
			v1.filestore.DeleteFile(img.MediaFile)
		}
		for _, vid := range videoFiles {
			v1.filestore.DeleteFile(vid.MediaFile)
		}
		if pr.MainImage != nil {
			v1.filestore.DeleteFile(*pr.MainImage)
		}
		v1.storage.Product().DeleteProductByID(context.Background(), pr.ID)
		return
	}

	if pr.MainImage != nil {
		pr.MainImage = models.GetStringAddress(v1.filestore.GetURL(*pr.MainImage))
	}
	pr.ImageFiles = imageFiles
	pr.VideoFiles = videoFiles

	v1.response(c, http.StatusOK, pr)
}

type GetAllProductsQueryParams struct {
	CategoryID *string `form:"cid"`
	Query      *string `form:"q"`
	BrandID    *string `form:"bid"`
	Offset     int     `form:"offset"`
	Limit      int     `form:"limit"`
}

// GetAllProducts
// @id getAllProducts
// @router /api/product [get]
// @summary get all products
// @description get all products
// @tags product
// @param cid query string false "Category ID to search in"
// @param q query string false "Query to search product"
// @param bid query string false "Brand ID to search in"
// @param offset query int false "Offset value. Default 0"
// @param limit query int false "Limit value. Default 10"
// @produce json
// @Success 200 {object} []models.Product "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid / status invalid"
// @Failure 404 {object} models_v1.Response "Category not found / Brand not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetAllProducts(c *gin.Context) {
	var m GetAllProductsQueryParams
	if err := c.ShouldBind(&m); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("bad request", logs.Error(err))
		return
	}
	if m.CategoryID != nil {
		if !helper.IsValidUUID(*m.CategoryID) {
			v1.error(c, status.StatusBadUUID)
			return
		}
	}
	if m.BrandID != nil {
		if !helper.IsValidUUID(*m.BrandID) {
			v1.error(c, status.StatusBadUUID)
			return
		}
	}
	products, err := v1.storage.Product().GetAll(context.Background(),
		m.Query,
		m.CategoryID,
		m.BrandID,
		models.GetProductAllLimits{
			Limit:  m.Limit,
			Offset: m.Offset,
		})
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get all products", logs.Error(err))
		return
	}

	v1.response(c, http.StatusOK, products)
}

// GetProductByID
// @id getProductByID
// @router /api/product/{id} [get]
// @summary get product by id
// @description get product by id
// @tags product
// @param id path string true "product id"
// @produce json
// @Success 200 {object} models.Product "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid"
// @Failure 404 {object} models_v1.Response "Product not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	product, err := v1.storage.Product().GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusProductNotFount)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not find product by id", logs.String("product_id", id),
			logs.Error(err))
		return
	}
	if product.DeletedAt != nil {
		v1.error(c, status.StatusProductNotFount)
		return
	}
	imgFiles, err := v1.storage.Product().GetProductImageFilesByID(context.Background(), id)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get image files for a product", logs.Error(err), logs.String("product_id", id))
	}
	product.ImageFiles = imgFiles

	vdFiles, err := v1.storage.Product().GetProductVideoFilesByID(context.Background(), id)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get video files for a product", logs.Error(err), logs.String("product_id", id))
	}
	product.VideoFiles = vdFiles

	v1.response(c, http.StatusOK, product)
}
