class CreateContactPostalAddresses < ActiveRecord::Migration[7.0]
  def change
    create_table :contact_postal_addresses, id: :string do |t|
      t.references :contact, null: false, type: :string, foreign_key: { on_delete: :cascade }
      t.string :name
      t.string :street
      t.string :locality
      t.string :region
      t.string :postal_code
      t.string :country
      t.boolean :preferred, null: false, default: false

      t.timestamps
    end
  end
end
