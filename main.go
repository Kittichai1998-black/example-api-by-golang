package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// โครงสร้างข้อมูล (Struct) สำหรับข้อมูลสินค้า
type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// เก็บรายการสินค้าในตัวแปรจำลอง (ในความจริงจะใช้ Database)
var products []Product

var books []Book

// ฟังก์ชัน GET: ดึงข้อมูลสินค้าทั้งหมด
func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// ฟังก์ชัน POST: เพิ่มสินค้าใหม่
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Product
	_ = json.NewDecoder(r.Body).Decode(&product)
	products = append(products, product)
	json.NewEncoder(w).Encode(product)
}

// ฟังก์ชัน PUT: อัปเดตข้อมูลสินค้า
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // ดึง parameter จาก URL
	for index, item := range products {
		if item.ID == params["id"] {
			products = append(products[:index], products[index+1:]...)
			var product Product
			_ = json.NewDecoder(r.Body).Decode(&product)
			product.ID = params["id"]
			products = append(products, product)
			json.NewEncoder(w).Encode(product)
			return
		}
	}
	json.NewEncoder(w).Encode(&Product{})
}

// ฟังก์ชัน DELETE: ลบสินค้า
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range products {
		if item.ID == params["id"] {
			products = append(products[:index], products[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(products)
}

// ฟังก์ชันสำหรับสร้าง Book (POST)
func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// ฟังก์ชันสำหรับดึงข้อมูล Book ทั้งหมด (GET)
func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

// ฟังก์ชันสำหรับดึง Book ตาม ID (GET)
func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

// ฟังก์ชันสำหรับอัปเดต Book (PUT)
func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, item := range books {
		if item.ID == params["id"] {
			books = append(books[:i], books[i+1:]...)
			var updatedBook Book
			json.NewDecoder(r.Body).Decode(&updatedBook)
			updatedBook.ID = params["id"]
			books = append(books, updatedBook)
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

// ฟังก์ชันสำหรับลบ Book (DELETE)
func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, item := range books {
		if item.ID == params["id"] {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

// ฟังก์ชันหลัก (Main)
func main() {
	// ใส่ข้อมูลตัวอย่าง
	products = append(products, Product{ID: "1", Name: "Apple", Price: 100})
	products = append(products, Product{ID: "2", Name: "Banana", Price: 50})

	// สร้าง router
	r := mux.NewRouter()

	// กำหนด route สำหรับ API แต่ละประเภท

	// Routing Products
	r.HandleFunc("/products", GetProducts).Methods("GET")
	r.HandleFunc("/products", CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", DeleteProduct).Methods("DELETE")

	// Routing Books
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// รันเซิร์ฟเวอร์บนพอร์ต 8000
	fmt.Println("Server is running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
