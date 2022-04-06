class CreateContactAddresses < ActiveRecord::Migration[7.0]
  def change
    create_table :contact_addresses, id: :string do |t|
      t.references :contact, null: false, type: :string, foreign_key: { on_delete: :cascade }
      t.string :kind, null: false
      t.string :data, null: false
      t.string :name
      t.boolean :preferred, null: false, default: false

      t.timestamps
    end
  end
end
