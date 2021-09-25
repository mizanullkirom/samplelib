package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/mizanullkirom/samplelib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	RedisDB *redis.Client
	MongoDB *mongo.Database
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	var item samplelib.Item
	json.NewDecoder(r.Body).Decode(&item)
	item.Id = uuid.New().String()
	_, err := samplelib.AddOne(h.MongoDB, &item)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
	}
	w.Write([]byte("Item created successfully"))
	w.WriteHeader(201)
}

// GetItems function using get all items
func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	items := samplelib.Get(h.MongoDB, bson.M{})
	json.NewEncoder(w).Encode(items)
}

//GetItem function using get one item by id
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	item := &samplelib.Item{}
	val, err := samplelib.GetItem(h.RedisDB, id)
	if err == redis.Nil {
		err := samplelib.GetOne(h.MongoDB, item, bson.M{"item": item})
		if err != nil {
			fmt.Println(err)
			return
		}
		err = samplelib.SetItem(h.RedisDB, id, item)
		json.NewEncoder(w).Encode(item)
	} else if err != nil {
		fmt.Println(err)
		return
	}
	w.Write([]byte(fmt.Sprintf(val)))
	w.WriteHeader(200)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	existingItem := &samplelib.Item{}
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	// val,err = samplelib.GetItem()
	err := samplelib.GetOne(h.MongoDB, existingItem, bson.M{"id": id})
	if err != nil {
		http.Error(w, fmt.Sprintf("Item doesn't exist"), 400)
		return
	}
	_, err = samplelib.RemoveOne(h.MongoDB, bson.M{"id": id})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
	w.Write([]byte("Item deleted"))
	w.WriteHeader(200)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	item := &samplelib.Item{}
	json.NewDecoder(r.Body).Decode(item)
	_, err := samplelib.Update(h.MongoDB, bson.M{"id": id}, item)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
	w.Write([]byte("Contact update successful"))
	w.WriteHeader(200)
}
