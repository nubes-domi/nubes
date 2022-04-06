# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[7.0].define(version: 2022_04_05_193911) do
  create_table "contact_addresses", id: :string, force: :cascade do |t|
    t.string "contact_id", null: false
    t.string "kind", null: false
    t.string "data", null: false
    t.string "name"
    t.boolean "preferred", default: false, null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["contact_id"], name: "index_contact_addresses_on_contact_id"
  end

  create_table "contact_postal_addresses", id: :string, force: :cascade do |t|
    t.string "contact_id", null: false
    t.string "name"
    t.string "street"
    t.string "locality"
    t.string "region"
    t.string "postal_code"
    t.string "country"
    t.boolean "preferred", default: false, null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["contact_id"], name: "index_contact_postal_addresses_on_contact_id"
  end

  create_table "contacts", id: :string, force: :cascade do |t|
    t.string "user_id", null: false
    t.string "title"
    t.string "given_name"
    t.string "middle_name"
    t.string "family_name"
    t.string "suffix"
    t.string "nickname"
    t.string "gender"
    t.string "pronouns"
    t.string "birthdate"
    t.boolean "birthdate_has_year"
    t.boolean "birthdate_age_based"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["user_id"], name: "index_contacts_on_user_id"
  end

  create_table "user_sessions", id: :string, force: :cascade do |t|
    t.string "user_id", null: false
    t.datetime "expires_at", null: false
    t.string "ip_address"
    t.string "user_agent"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["user_id"], name: "index_user_sessions_on_user_id"
  end

  create_table "users", id: :string, force: :cascade do |t|
    t.string "username", null: false
    t.string "password_digest"
    t.boolean "admin"
    t.string "given_name"
    t.string "middle_name"
    t.string "family_name"
    t.string "email"
    t.boolean "email_verified"
    t.string "phone_number"
    t.boolean "phone_number_verified"
    t.date "birthdate"
    t.string "gender"
    t.string "pronouns"
    t.string "address_street"
    t.string "address_locality"
    t.string "address_region"
    t.string "address_postal_code"
    t.string "address_country"
    t.string "locale"
    t.string "zoneinfo"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["username"], name: "index_users_on_username", unique: true
  end

  add_foreign_key "contact_addresses", "contacts", on_delete: :cascade
  add_foreign_key "contact_postal_addresses", "contacts", on_delete: :cascade
  add_foreign_key "contacts", "users", on_delete: :cascade
  add_foreign_key "user_sessions", "users", on_delete: :cascade
end
