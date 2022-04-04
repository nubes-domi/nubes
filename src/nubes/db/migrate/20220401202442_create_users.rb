class CreateUsers < ActiveRecord::Migration[7.0]
  def change
    create_table :users, id: :string do |t|
      t.string :username, null: false, unique: true
      t.string :password_digest
      t.boolean :admin
      t.string :given_name
      t.string :middle_name
      t.string :family_name
      t.string :email
      t.boolean :email_verified
      t.string :phone_number
      t.boolean :phone_number_verified
      t.date :birthdate
      t.string :gender
      t.string :pronouns
      t.string :address_street
      t.string :address_locality
      t.string :address_region
      t.string :address_postal_code
      t.string :address_country
      t.string :locale
      t.string :zoneinfo

      t.timestamps
    end

    add_index :users, :username, unique: true
  end
end
