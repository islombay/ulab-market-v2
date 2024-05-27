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
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// CreateProduct
// @ID createProduct
// @Router /api/product [post]
// @Summary create product
// @Description create product
// @security ApiKeyAuth
// @accept multipart/form-data
// @tags product
// @produce json
// @param createProduct formData models_v1.CreateProduct true "create product"
// @param main_image formData file true "Main image"
// @param image_files formData []file false "Image files (multiple)"
// @param video_files formData []file false "Video files (multiple)"
// @Success 200 {object} models.Product "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid / status invalid"
// @Failure 404 {object} models_v1.Response "Category not found / Brand not found"
// @failure 409 {object} models_v1.Response "Articul already found"
// @Failure 413 {object} models_v1.Response "Image size is big / Video size is big / Image count too many / Video count too many / Articul too long"
// @Failure 415 {object} models_v1.Response "Image type is not supported / Video type is not supported"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) CreateProduct(c *gin.Context) {
	var m models_v1.CreateProduct
	if err := c.Bind(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		v1.log.Error("bad request", logs.Error(err))
		return
	}

	if m.IncomePrice <= 0 {
		v1.error(c, status.StatusProductPriceInvalid)
		return
	}

	if !helper.IsValidUUID(m.BranchID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	if _, err := v1.storage.Branch().GetByID(context.Background(), m.BranchID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusBranchNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get branch by id", logs.Error(err),
			logs.String("bid", m.BranchID))
		return
	}

	if len(m.Articul) > 250 {
		v1.error(c, status.StatusProductArticulTooLong)
		return
	}

	// check whether the same articul exists
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

	if m.OutcomePrice <= 0 {
		v1.error(c, status.StatusProductPriceInvalid)
		return
	}

	// is valid status provided
	m.Status = strings.ToLower(m.Status)
	if m.Status == "" {
		m.Status = helper.ProductStatusActive
	}
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
		OutcomePrice:  m.OutcomePrice,
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
			v1.log.Error("could not check file for validity", logs.Error(err))
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
			v1.log.Error("could not check file for validity", logs.Error(err))
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
			v1.log.Error("could not check file for validity", logs.Error(err))
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

	incomeCreate := models_v1.CreateIncome{
		BranchID: m.BranchID,
		Comment:  "created during product create",
		// TotalPrice: m.IncomePrice * float32(m.Quantity),
		Products: []models_v1.CreateIncomeProduct{
			models_v1.CreateIncomeProduct{
				ProductID:    pr.ID,
				Quantity:     int(m.Quantity),
				ProductPrice: m.IncomePrice,
			},
		},
	}

	v1.service.Income().Create(context.Background(), incomeCreate)
	v1.response(c, http.StatusOK, pr)
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
// @param page query int false "Page value. Default 0"
// @param limit query int false "Limit value. Default 10"
// @produce json
// @Success 200 {object} []models.Product "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid / status invalid"
// @Failure 404 {object} models_v1.Response "Category not found / Brand not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetAllProducts(c *gin.Context) {
	var m models_v1.GetAllProductsQueryParams
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
			Offset: (m.Page - 1) * m.Limit,
		})
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get all products", logs.Error(err))
		return
	}

	res := make([]models_v1.Product, len(products))

	for i, p := range products {
		if p.MainImage != nil {
			p.MainImage = models.GetStringAddress(v1.filestore.GetURL(*p.MainImage))
		}
		tmp := models_v1.Product{}
		if err := helper.Reobject(*p, &tmp, "obj"); err != nil {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not reobject", logs.Error(err))
			return
		}

		if p.CategoryID != nil {
			tmpCat, err := v1.storage.Category().GetByID(context.Background(), *p.CategoryID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					v1.log.Error("could not find category by id from product", logs.String("pid", p.ID),
						logs.String("cid", *p.CategoryID))
				} else {
					v1.log.Error("could not get category information from product", logs.String("pid", p.ID),
						logs.String("cid", *p.CategoryID), logs.Error(err))
				}
			} else {
				tmp.CategoryInformation = *tmpCat
				if tmp.CategoryInformation.Image != nil {
					tmp.CategoryInformation.Image = models.GetStringAddress(v1.filestore.GetURL(*tmp.CategoryInformation.Image))
				}
			}
		}
		tmp.Articul = p.Articul

		res[i] = tmp
	}

	v1.response(c, http.StatusOK, res)
}

// GetAllProductsPagination
// @id GetAllProductsPagination
// @router /api/product/_pagin [get]
// @summary get all products
// @description get all products
// @tags product
// @param cid query string false "Category ID to search in"
// @param q query string false "Query to search product"
// @param bid query string false "Brand ID to search in"
// @param page query int false "page value. Default 1"
// @param limit query int false "Limit value. Default 10"
// @produce json
// @Success 200 {object} []models.Product "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid / status invalid"
// @Failure 404 {object} models_v1.Response "Category not found / Brand not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetAllProductsPagination(c *gin.Context) {
	var m models_v1.ProductPagination
	c.ShouldBind(&m)

	m.Fix()

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
	products, count, err := v1.storage.Product().GetAllPagination(context.Background(), m)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get all products", logs.Error(err))
		return
	}

	res := make([]models_v1.Product, len(products))

	for i, p := range products {
		if p.MainImage != nil {
			p.MainImage = models.GetStringAddress(v1.filestore.GetURL(*p.MainImage))
		}
		tmp := models_v1.Product{}
		if err := helper.Reobject(*p, &tmp, "obj"); err != nil {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not reobject", logs.Error(err))
			return
		}

		if p.CategoryID != nil {
			tmpCat, err := v1.storage.Category().GetByID(context.Background(), *p.CategoryID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					v1.log.Error("could not find category by id from product", logs.String("pid", p.ID),
						logs.String("cid", *p.CategoryID))
				} else {
					v1.log.Error("could not get category information from product", logs.String("pid", p.ID),
						logs.String("cid", *p.CategoryID), logs.Error(err))
				}
			} else {
				tmp.CategoryInformation = *tmpCat
				if tmp.CategoryInformation.Image != nil {
					tmp.CategoryInformation.Image = models.GetStringAddress(v1.filestore.GetURL(*tmp.CategoryInformation.Image))
				}
			}
		}
		tmp.Articul = p.Articul

		res[i] = tmp
	}

	v1.response(c, http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Count:      count,
		Data:       res,
	})
}

// GetProductByID
// @id getProductByID
// @router /api/product/{id} [get]
// @summary get product by id
// @description get product by id
// @tags product
// @param id path string true "product id"
// @produce json
// @Success 200 {object} models_v1.Product "Success"
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

	var tmp models_v1.Product

	imgFiles, err := v1.storage.Product().GetProductImageFilesByID(context.Background(), id)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get image files for a product", logs.Error(err), logs.String("product_id", id))
	}

	tmp.ImageFiles = make([]models_v1.ProductMediaFiles, len(imgFiles))
	for i, _ := range imgFiles {
		imgFiles[i].MediaFile = v1.filestore.GetURL(imgFiles[i].MediaFile)
		tmp.ImageFiles[i] = models_v1.ProductMediaFiles{
			ID:        imgFiles[i].ID,
			MediaFile: imgFiles[i].MediaFile,
		}
	}

	vdFiles, err := v1.storage.Product().GetProductVideoFilesByID(context.Background(), id)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get video files for a product", logs.Error(err), logs.String("product_id", id))
	}

	tmp.VideoFiles = make([]models_v1.ProductMediaFiles, len(vdFiles))
	for i, _ := range vdFiles {
		vdFiles[i].MediaFile = v1.filestore.GetURL(vdFiles[i].MediaFile)
		tmp.VideoFiles[i] = models_v1.ProductMediaFiles{
			ID:        vdFiles[i].ID,
			MediaFile: vdFiles[i].MediaFile,
		}
	}

	if product.MainImage != nil {
		product.MainImage = models.GetStringAddress(v1.filestore.GetURL(*product.MainImage))
	}

	if err := helper.Reobject(*product, &tmp, "obj"); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not reobject", logs.Error(err))
		return
	}

	if product.MainImage != nil {
		tmp.MainImage = *product.MainImage
	}

	if err := v1.storage.Product().IncrementViewCount(context.Background(), tmp.ID); err != nil {
		v1.log.Error("could not increment view count for product", logs.Error(err),
			logs.String("product_id", tmp.ID))
	}

	v1.response(c, http.StatusOK, tmp)
}

// GetProductByIDAdmin
// @id GetProductByIDAdmin
// @router /api/product/admin/{id} [get]
// @summary get product by id
// @description get product by id
// @tags product
// @security ApiKeyAuth
// @param id path string true "product id"
// @produce json
// @Success 200 {object} models_v1.Product "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid"
// @Failure 404 {object} models_v1.Response "Product not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetProductByIDAdmin(c *gin.Context) {
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

	var tmp models_v1.Product

	imgFiles, err := v1.storage.Product().GetProductImageFilesByID(context.Background(), id)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get image files for a product", logs.Error(err), logs.String("product_id", id))
	}

	tmp.ImageFiles = make([]models_v1.ProductMediaFiles, len(imgFiles))
	for i, _ := range imgFiles {
		imgFiles[i].MediaFile = v1.filestore.GetURL(imgFiles[i].MediaFile)
		tmp.ImageFiles[i] = models_v1.ProductMediaFiles{
			ID:        imgFiles[i].ID,
			MediaFile: imgFiles[i].MediaFile,
		}
	}

	vdFiles, err := v1.storage.Product().GetProductVideoFilesByID(context.Background(), id)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get video files for a product", logs.Error(err), logs.String("product_id", id))
	}

	tmp.VideoFiles = make([]models_v1.ProductMediaFiles, len(vdFiles))
	for i, _ := range vdFiles {
		vdFiles[i].MediaFile = v1.filestore.GetURL(vdFiles[i].MediaFile)
		tmp.VideoFiles[i] = models_v1.ProductMediaFiles{
			ID:        vdFiles[i].ID,
			MediaFile: vdFiles[i].MediaFile,
		}
	}

	if product.MainImage != nil {
		product.MainImage = models.GetStringAddress(v1.filestore.GetURL(*product.MainImage))
	}

	if err := helper.Reobject(*product, &tmp, "obj"); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not reobject", logs.Error(err))
		return
	}
	tmp.Articul = product.Articul

	if product.MainImage != nil {
		tmp.MainImage = *product.MainImage
	}

	fmt.Println(tmp.CategoryID)
	if tmp.CategoryID != "" {
		tmpCat, err := v1.storage.Category().GetByID(context.Background(), tmp.CategoryID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.log.Error("could not find category by id from product", logs.String("pid", tmp.ID),
					logs.String("cid", tmp.CategoryID))
			} else {
				v1.log.Error("could not get category information from product", logs.String("pid", tmp.ID),
					logs.String("cid", tmp.CategoryID), logs.Error(err))
			}
		} else {
			tmp.CategoryInformation = *tmpCat
			if tmp.CategoryInformation.Image != nil {
				tmp.CategoryInformation.Image = models.GetStringAddress(v1.filestore.GetURL(*tmp.CategoryInformation.Image))
			}
		}
	}
	v1.response(c, http.StatusOK, tmp)
}

// ChangeProductMainImage
// @id ChangeProductMainImage
// @router /api/product/change_image [post]
// @summary change image of a product
// @description change image of a product
// @tags product
// @security ApiKeyAuth
// @param changeProductMainImage formData models_v1.ChangeProductMainImage true "body"
// @param image formData file true "image"
// @produce json
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid"
// @Failure 404 {object} models_v1.Response "Product not found"
// @failure 413 {object} models_v1.Response "Main image size exceeds the limit"
// @failure 415 {object} models_v1.Response "Unsupported media type"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) ChangeProductMainImage(c *gin.Context) {
	var m models_v1.ChangeProductMainImage
	if err := c.Bind(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}
	if !helper.IsValidUUID(m.ProductID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	product, err := v1.storage.Product().GetByID(context.Background(), m.ProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusProductNotFount)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get product by id", logs.Error(err),
			logs.String("product_id", m.ProductID))
		return
	}

	if st := v1.ValidateImage(m.Image); st != nil {
		fmt.Println(st.Code)
		v1.error(c, *st)
		return
	}

	now := time.Now()
	imageURL := product.ID + "-" + strconv.FormatInt(now.Unix(), 10)
	internalURL, err := v1.filestore.Create(m.Image, filestore.FolderProduct, imageURL)
	if err != nil {
		v1.error(c, status.StatusInternal)
		return
	}

	if err := v1.storage.Product().ChangeMainImage(context.Background(), product.ID, internalURL, now); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not change main image", logs.Error(err),
			logs.String("product_id", m.ProductID))
		return
	}

	if product.MainImage != nil {
		defer func() {
			if err := v1.filestore.DeleteFile(*product.MainImage); err != nil {
				v1.log.Debug("could not delete image file", logs.Error(err))
			}
		}()
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

func (v1 *Handlers) ValidateImage(img *multipart.FileHeader) *status.Status {
	if _, err, msg := helper.IsValidImage(img); err != nil {
		if errors.Is(err, helper.ErrInvalidImageType) {
			v1.log.Error("got invalid image extension", logs.String("got", msg))
			return &status.StatusImageTypeUnkown
		}
		v1.log.Error("could not check file for validity", logs.Error(err))
		return &status.StatusInternal
	}
	if img.Size > v1.cfg.Media.ProductPhotoMaxSize {
		return &status.StatusProductMainImageMaxSizeExceed
	}
	return nil
}

// DeleteProduct
// @id DeleteProduct
// @router /api/product/{id} [delete]
// @summary delete product
// @description delete product
// @tags product
// @security ApiKeyAuth
// @param id path string true "product id"
// @produce json
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad ID"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	if err := v1.storage.Product().DeleteProductByID(context.Background(), id); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not delete product", logs.Error(err), logs.String("product_id", id))
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// AddProductImageFiles
// @id AddProductImageFiles
// @router /api/product/add_image_files [post]
// @summary add image files to product
// @description add image files to product
// @security ApiKeyAuth
// @accept multipart/form-data
// @tags product
// @produce json
// @param addImageFiles formData models_v1.AddProductMediaFiles true "add product image files"
// @param media_files formData file true "image files"
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid"
// @Failure 404 {object} models_v1.Response "Product not found"
// @Failure 413 {object} models_v1.Response "Image size is big / Image count too many"
// @Failure 415 {object} models_v1.Response "Image type is not supported"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) AddProductImageFiles(c *gin.Context) {
	var m models_v1.AddProductMediaFiles
	if err := c.Bind(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		v1.log.Debug("invalid binding request add image files struct", logs.Error(err))
		return
	}

	if !helper.IsValidUUID(m.ProductID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	if _, err := v1.storage.Product().GetByID(context.Background(), m.ProductID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusProductNotFount)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not find product by id", logs.Error(err),
			logs.String("product_id", m.ProductID),
		)
		return
	}

	imgFiles, err := v1.storage.Product().GetProductImageFilesByID(context.Background(), m.ProductID)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get image files for a product", logs.Error(err), logs.String("product_id", m.ProductID))
		return
	}

	if len(imgFiles)+1 > v1.cfg.Media.ProductPhotoMaxCount {
		v1.error(c, status.StatusProductPhotoMaxCount)
		return
	}
	// var imageFilesStopped struct {
	// 	Status *status.Status
	// }

	// imageFiles := make([]models.ProductMediaFiles, 1)

	if _, err, msg := helper.IsValidImage(m.MediaFiles); err != nil {
		if errors.Is(err, helper.ErrInvalidImageType) {
			v1.error(c, status.StatusImageTypeUnkown)
			v1.log.Error("got invalid image extension", logs.String("got", msg))
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not check file for validity", logs.Error(err))
		return
	}
	if m.MediaFiles.Size > v1.cfg.Media.ProductPhotoMaxSize {
		v1.error(c, status.StatusImageMaxSizeExceed)
		return
	}
	id := uuid.NewString()
	url, err := v1.filestore.Create(m.MediaFiles, filestore.FolderProduct, id)
	if err != nil {
		v1.error(c, status.StatusInternal)
		return
	}
	if err := v1.storage.Product().CreateProductImageFile(context.Background(),
		id, m.ProductID, url,
	); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not create product image file in db", logs.Error(err),
			logs.String("pid", m.ProductID),
		)
		return
	}
	// for ind, i := range m.MediaFiles {
	// 	if _, err, msg := helper.IsValidImage(i); err != nil {
	// 		if errors.Is(err, helper.ErrInvalidImageType) {
	// 			imageFilesStopped.Status = &status.StatusImageTypeUnkown
	// 			v1.log.Error("got invalid image extension", logs.String("got", msg))
	// 			break
	// 		}
	// 		imageFilesStopped.Status = &status.StatusInternal
	// 		v1.log.Error("could not check file for validity", logs.Error(err))
	// 		break
	// 	}
	// 	if i.Size > v1.cfg.Media.ProductPhotoMaxSize {
	// 		imageFilesStopped.Status = &status.StatusImageMaxSizeExceed
	// 		break
	// 	}
	// 	id := uuid.NewString()
	// 	url, err := v1.filestore.Create(i, filestore.FolderProduct, id)
	// 	if err != nil {
	// 		imageFilesStopped.Status = &status.StatusInternal
	// 		break
	// 	}
	// 	if err := v1.storage.Product().CreateProductImageFile(context.Background(),
	// 		id, m.ProductID, url,
	// 	); err != nil {
	// 		imageFilesStopped.Status = &status.StatusInternal
	// 		v1.log.Error("could not create product image file in db", logs.Error(err),
	// 			logs.String("pid", m.ProductID),
	// 		)
	// 		break
	// 	}
	// 	imageFiles[ind] = models.ProductMediaFiles{
	// 		ID:        id,
	// 		ProductID: m.ProductID,
	// 		MediaFile: v1.filestore.GetURL(url),
	// 	}
	// }

	// if imageFilesStopped.Status != nil {
	// 	v1.error(c, *imageFilesStopped.Status)
	// 	for _, img := range imageFiles {
	// 		v1.filestore.DeleteFile(img.MediaFile)
	// 	}
	// 	return
	// }

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// AddProductVideoFiles
// @id AddProductVideoFiles
// @router /api/product/add_video_files [post]
// @summary add video files to product
// @description add video files to product
// @security ApiKeyAuth
// @accept multipart/form-data
// @tags product
// @produce json
// @param addImageFiles formData models_v1.AddProductMediaFiles true "add product image files"
// @param media_files formData []file true "video files"
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid"
// @Failure 404 {object} models_v1.Response "Product not found"
// @Failure 413 {object} models_v1.Response "Video size is big / Video count too many"
// @Failure 415 {object} models_v1.Response "Video type is not supported"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) AddProductVideoFiles(c *gin.Context) {
	var m models_v1.AddProductVideoFiles
	if err := c.Bind(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		v1.log.Debug("invalid binding request add video files struct", logs.Error(err))
		return
	}

	if !helper.IsValidUUID(m.ProductID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	if _, err := v1.storage.Product().GetByID(context.Background(), m.ProductID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusProductNotFount)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not find product by id", logs.Error(err),
			logs.String("product_id", m.ProductID),
		)
		return
	}

	if len(m.MediaFiles) > v1.cfg.Media.ProductVideoMaxCount {
		v1.error(c, status.StatusProductVideoMaxCount)
		return
	}
	var imageFilesStopped struct {
		Status *status.Status
	}

	imageFiles := make([]models.ProductMediaFiles, len(m.MediaFiles))
	for ind, i := range m.MediaFiles {
		if _, err, msg := helper.IsValidImage(i); err != nil {
			if errors.Is(err, helper.ErrInvalidVideoType) {
				imageFilesStopped.Status = &status.StatusVideoTypeUnkown
				v1.log.Error("got invalid video extension", logs.String("got", msg))
				break
			}
			imageFilesStopped.Status = &status.StatusInternal
			v1.log.Error("could not check file for validity", logs.Error(err))
			break
		}
		if i.Size > v1.cfg.Media.ProductVideoMaxSize {
			imageFilesStopped.Status = &status.StatusProductVideoMaxSizeExceed
			break
		}
		id := uuid.NewString()
		url, err := v1.filestore.Create(i, filestore.FolderProduct, id)
		if err != nil {
			imageFilesStopped.Status = &status.StatusInternal
			break
		}
		if err := v1.storage.Product().CreateProductVideoFile(context.Background(),
			id, m.ProductID, url,
		); err != nil {
			imageFilesStopped.Status = &status.StatusInternal
			v1.log.Error("could not create product video file in db", logs.Error(err),
				logs.String("pid", m.ProductID),
			)
			break
		}
		imageFiles[ind] = models.ProductMediaFiles{
			ID:        id,
			ProductID: m.ProductID,
			MediaFile: v1.filestore.GetURL(url),
		}
	}
	if imageFilesStopped.Status != nil {
		v1.error(c, *imageFilesStopped.Status)
		for _, img := range imageFiles {
			v1.filestore.DeleteFile(img.MediaFile)
		}
		return
	}
	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// ChangeProductPrice
// @id ChangeProductPrice
// @router /api/product/change_price [put]
// @tags product
// @security ApiKeyAuth
// @Summary Change products price (selling price)
// @description change products price
// @accept json
// @produce json
// @param price_body body models_v1.ChangeProductPrice true "Change product price"
// @success 200 {object} models_v1.Response "Success"
// @failure 400 {object} models_v1.Response "Bad id/ bad price"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) ChangeProductPrice(c *gin.Context) {
	var m models_v1.ChangeProductPrice
	if err := c.BindJSON(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}
	if !helper.IsValidUUID(m.ID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	if m.Price <= 0 {
		v1.error(c, status.StatusProductPriceInvalid)
		return
	}

	err := v1.storage.Product().ChangeProductPrice(context.Background(), m.ID, m.Price)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not update price of product", logs.Error(err),
			logs.String("pid", m.ID), logs.Any("price", m.Price))
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Message: "Ok",
		Code:    200,
	})
}

// EditProduct
// @id			EditProduct
// @router		/api/product [put]
// @summary		edit product
// @description	edit product
// @tags		product
// @accept		json
// @produce		json
// @security	ApiKeyAuth
// @param		edit_product body models_v1.ChangeProductRequest true "edit product request"
// @success		200 {object} models.Product "Successfull edit"
// @failure		400 {object} models_v1.Response "Bad request/bad id/ bad price/ bad status"
// @failure		404 {object} models_v1.Response "Product not found/ Category not found/ Brand not found"
// @failure		500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) EditProduct(c *gin.Context) {
	var m models_v1.ChangeProductRequest
	if err := c.BindJSON(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		v1.log.Error("bad request", logs.Error(err))
		return
	}

	if !helper.IsValidUUID(m.ID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	_, err := v1.storage.Product().GetByID(context.Background(), m.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusProductNotFount)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not find product by id", logs.Error(err),
			logs.String("pid", m.ID))
		return
	}

	if m.OutcomePrice <= 0 {
		v1.error(c, status.StatusProductPriceInvalid)
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

	productModel := models.Product{
		ID:      m.ID,
		Articul: m.Articul,

		NameRu: m.NameRu,
		NameUz: m.NameUz,

		DescriptionUz: m.DescriptionUz,
		DescriptionRu: m.DescriptionRu,

		OutcomePrice: m.OutcomePrice,
		CategoryID:   &m.CategoryID,
		BrandID:      &m.BrandID,

		Status: m.Status,
	}

	v1.storage.Product().Change(context.Background(), &productModel)
	if productModel.MainImage != nil {
		productModel.MainImage = models.GetStringAddress(v1.filestore.GetURL(*productModel.MainImage))
	}
	v1.response(c, http.StatusOK, productModel)
}
