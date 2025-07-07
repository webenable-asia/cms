#!/bin/bash

# MongoDB initialization script
# This script creates the CMS database and a user with appropriate permissions

echo "Creating CMS database and user..."

mongosh <<EOF
use admin

// Create CMS database
use cms

// Create application user
db.createUser({
  user: "${MONGO_CMS_USER}",
  pwd: "${MONGO_CMS_PASSWORD}",
  roles: [
    {
      role: "readWrite",
      db: "cms"
    }
  ]
})

// Create initial collections with validation
db.createCollection("posts", {
  validator: {
    \$jsonSchema: {
      bsonType: "object",
      required: ["title", "content", "author", "created_at"],
      properties: {
        title: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        content: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        author: {
          bsonType: "objectId",
          description: "must be an objectId and is required"
        },
        created_at: {
          bsonType: "date",
          description: "must be a date and is required"
        },
        updated_at: {
          bsonType: "date",
          description: "must be a date"
        },
        tags: {
          bsonType: "array",
          items: {
            bsonType: "string"
          }
        },
        status: {
          enum: ["draft", "published", "archived"],
          description: "can only be one of the enum values"
        }
      }
    }
  }
})

db.createCollection("users", {
  validator: {
    \$jsonSchema: {
      bsonType: "object",
      required: ["username", "email", "created_at"],
      properties: {
        username: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        email: {
          bsonType: "string",
          pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\$",
          description: "must be a valid email and is required"
        },
        password_hash: {
          bsonType: "string",
          description: "must be a string"
        },
        created_at: {
          bsonType: "date",
          description: "must be a date and is required"
        },
        updated_at: {
          bsonType: "date",
          description: "must be a date"
        }
      }
    }
  }
})

// Create indexes for better performance
db.posts.createIndex({ "title": "text", "content": "text" })
db.posts.createIndex({ "author": 1 })
db.posts.createIndex({ "created_at": -1 })
db.posts.createIndex({ "status": 1 })
db.posts.createIndex({ "tags": 1 })

db.users.createIndex({ "username": 1 }, { unique: true })
db.users.createIndex({ "email": 1 }, { unique: true })

print("CMS database initialization completed successfully!")
EOF
