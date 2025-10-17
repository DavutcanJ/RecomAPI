package repository

import (
    "context"
    "errors"
    "time"

    "github.com/DavutcanJ/customer-segmentation-api/internal/config"
    "github.com/DavutcanJ/customer-segmentation-api/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        collection: config.DB.Collection("users"),
    }
}

func (r *UserRepository) CreateUser(user *models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    user.ID = primitive.NewObjectID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    user.IsActive = true
    user.Role = "user" // Default role

    _, err := r.collection.InsertOne(ctx, user)
    return err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid user ID")
    }

    var user models.User
    err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) CreateIndexes() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Create unique indexes for performance and data integrity
    indexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "username", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys: bson.D{{Key: "created_at", Value: 1}},
        },
    }

    _, err := r.collection.Indexes().CreateMany(ctx, indexes)
    return err
}