package articlecontroller

import (
	"errors"
	"go-api-article/config"
	"go-api-article/helper"
	"go-api-article/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var articles []models.Article
	var articlesResponse []models.ArticleResponse

	if err := config.DB.Joins("Category").Find(&articles).Find(&articlesResponse).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "List Artikel", articlesResponse)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var article models.Article

	// parse multipart/form-data request
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.Response(w, 400, "File terlalu besar", nil)
		return
	}

	// get form data
	title := r.FormValue("title")
	if title == "" {
		helper.Response(w, 400, "Title tidak boleh kosong", nil)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		helper.Response(w, 400, "Invalid category ID", nil)
		return
	}
	content := r.FormValue("content")
	if content == "" {
		helper.Response(w, 400, "Content tidak boleh kosong", nil)
		return
	}

	// get file from form data
	file, handler, err := r.FormFile("thumbnail")
	if err != nil {
		helper.Response(w, 400, "File thumbnail tidak ditemukan", nil)
		return
	}
	defer file.Close()

	// generate unique file name and save to server
	fileName := uuid.New().String() + filepath.Ext(handler.Filename)
	filePath := "thumbnails/" + fileName
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		helper.Response(w, 500, "Gagal menyimpan file thumbnail", nil)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	article.Title = title
	article.CategoryID = uint(categoryID)
	article.Content = content
	article.Thumbnail = fileName

	//check category
	var category models.Category
	if err := config.DB.First(&category, article.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Kategori Tidak Ditemukan", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	if err := config.DB.Create(&article).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Berhasil Menambahkan Artikel", nil)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var articles models.Article
	var articlesResponse models.ArticleResponse

	if err := config.DB.Joins("Category").First(&articles, id).First(&articlesResponse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Artikel Tidak Ditemukan", nil)
			return
		}
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "Detail Artikel", articlesResponse)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var article models.Article

	// get article ID from URL parameter
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		helper.Response(w, http.StatusBadRequest, "Invalid article ID", nil)
		return
	}

	// get article data from database
	if err := config.DB.First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Article Tidak Ditemukan", nil)
			return
		}
		helper.Response(w, 500, "Gagal Mendapatkan Data Artikel", nil)
		return
	}

	// parse form data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.Response(w, 400, "Ukuran File Terlalu Besar", nil)
		return
	}

	// update article data
	if title := r.FormValue("title"); title != "" {
		article.Title = title
	}

	if categoryID, err := strconv.Atoi(r.FormValue("category_id")); err == nil && categoryID != 0 {
		// check if category exists
		var category models.Category
		if err := config.DB.First(&category, categoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				helper.Response(w, 400, "Kategori Tidak Ditemukan", nil)
				return
			}
			helper.Response(w, 500, "Gagal Mendapatkan Data Kategori", nil)
			return
		}
		article.CategoryID = uint(categoryID)
	}

	if content := r.FormValue("content"); content != "" {
		article.Content = content
	}

	if file, handler, err := r.FormFile("thumbnail"); err == nil && file != nil {
		defer file.Close()

		// check if thumbnail already exists
		if article.Thumbnail != "" {
			if err := os.Remove("thumbnails/" + article.Thumbnail); err != nil {
				helper.Response(w, 500, "Failed to delete old thumbnail", nil)
				return
			}
		}

		// generate new file name
		ext := filepath.Ext(handler.Filename)
		fileName := uuid.New().String() + ext

		// save new thumbnail to server
		filePath := "thumbnails/" + fileName
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			helper.Response(w, 500, "Failed to save thumbnail file", nil)
			return
		}
		defer f.Close()

		io.Copy(f, file)
		article.Thumbnail = fileName
	}

	// save updated article data to database
	if err := config.DB.Save(&article).Error; err != nil {
		helper.Response(w, 500, "Failed to save updated article data", nil)
		return
	}

	helper.Response(w, 200, "Artikel Berhasil Diperbaharui", nil)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParams)
	if err != nil {
		helper.Response(w, 400, "Invalid article ID", nil)
		return
	}

	var article models.Article

	// Find article record by ID
	if err := config.DB.First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Artikel Tidak Ditemukan", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	// Delete article record from database
	result := config.DB.Delete(&article, id)
	if result.Error != nil {
		helper.Response(w, 500, result.Error.Error(), nil)
		return
	}

	// Delete old thumbnail file if it exists
	if article.Thumbnail != "" {
		if err := os.Remove("thumbnails/" + article.Thumbnail); err != nil {
			helper.Response(w, 500, "Gagal menghapus thumbnail lama", nil)
			return
		}
	}

	// Check if article record was deleted
	if result.RowsAffected == 0 {
		helper.Response(w, 404, "Artikel Tidak Ditemukan", nil)
		return
	}

	helper.Response(w, 200, "Artikel Berhasil Dihapus", nil)

}
