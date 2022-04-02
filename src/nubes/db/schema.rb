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

ActiveRecord::Schema[7.0].define(version: 2022_04_01_210232) do
  create_table "user_sessions", id: :string, force: :cascade do |t|
    t.integer "user_id"
    t.integer "{:type=>:string, :on_delete=>:cascade}_id"
    t.datetime "expires_at", null: false
    t.string "ip_address"
    t.string "user_agent"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["user_id"], name: "index_user_sessions_on_user_id"
    t.index ["{:type=>:string, :on_delete=>:cascade}_id"], name: "index_user_sessions_on_{:type=>:string, :on_delete=>:cascade}_id"
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

end
