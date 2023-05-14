package categorycontroller

import (
	"errors"
	"go-api-article/config"
	"go-api-article/helper"
	"go-api-article/models"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var category []models.Category

	if err := config.DB.Find(&category).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "List Category", category)

}

func Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category

	//get form data
	name := r.FormValue("name")
	if name == "" {
		helper.Response(w, 400, "Nama Kategori Tidak Boleh Kosong", nil)
		return
	}

	description := r.FormValue("description")
	if description == "" {
		helper.Response(w, 400, "Deskripsi Tidak Boleh Kosong", nil)
		return
	}

	defer r.Body.Close()

	category.Name = name
	category.Description = description

	if err := config.DB.Create(&category).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Berhasil Menambahkan Category", category)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var category models.Category

	if err := config.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Category Tidak Ditemukan", nil)
			return
		}
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "Detail Category", category)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var category models.Category

	idParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParams)

	if err != nil {
		helper.Response(w, 400, "Invalid Article ID", nil)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	defer r.Body.Close()

	if name == "" || description == "" {
		helper.Response(w, 400, "Name dan Description harus diisi", nil)
		return
	}

	//update category
	if err := config.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Kategori Tidak Ditemukan", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	category.Name = name
	category.Description = description

	if err := config.DB.Save(&category).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Berhasil Memperbaharui Category", category)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var category models.Category

	result := config.DB.Delete(&category, id)

	if result.Error != nil {
		helper.Response(w, 500, result.Error.Error(), nil)
		return
	}

	if result.RowsAffected == 0 {
		helper.Response(w, 404, "Category Tidak Ditemukan", nil)
		return
	}

	helper.Response(w, 200, "Category Berhasil Dihapus", nil)
}
