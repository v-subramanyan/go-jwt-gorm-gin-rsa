package models

import "time"

// User represents a system user with associated JWT tokens, RSA key, groups, and roles.
type User struct {
	ID           uint       `gorm:"primaryKey"`
	Name         string     `gorm:"unique;not null"`
	Email        string     `gorm:"unique;not null"`
	Password     string     `gorm:"not null"`
	JWTToken     string     // JWT Token for the user
	RefreshToken string     // Refresh token for the user
	RSAKey       RSAKeyPair `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // One-to-one relationship with RSAKey
	Groups       []Group    `gorm:"many2many:user_groups;"`                        // Many-to-many relationship with groups
	Roles        []Role     `gorm:"many2many:user_roles;"`                         // Many-to-many relationship with roles
}

// Group represents a group that a user can belong to, which can also have a parent group.
type Group struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;not null"`        // Unique name for the group
	Members  []User `gorm:"many2many:user_groups;"` // Many-to-many relationship with users
	Parent   *Group // Pointer to the parent group
	ParentID *uint  // Foreign key for the parent group (nullable)
}

// Role represents a role that a user can have, with possible associations to groups.
type Role struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  `gorm:"unique;not null"`        // Unique name for the role
	Users  []User  `gorm:"many2many:user_roles;"`  // Many-to-many relationship with users
	Groups []Group `gorm:"many2many:role_groups;"` // Many-to-many relationship with groups
}

// RSAKey represents the RSA public and private keys associated with a user.
type RSAKeyPair struct {
	ID         uint      `gorm:"primaryKey"`
	PrivateKey string    // RSA private key in PEM format
	PublicKey  string    // RSA public key in PEM format
	UserID     uint      `gorm:"unique"` // Foreign key to the User
	CreatedAt  time.Time // Time when the key was created
	ExpiresAt  time.Time // Expiration time of the key
	IsActive   bool      // Whether the key is active or not
}
