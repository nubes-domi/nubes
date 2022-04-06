class CreateContacts < ActiveRecord::Migration[7.0]
  def change
    create_table :contacts, id: :string do |t|
      t.references :user, null: false, type: :string, foreign_key: { on_delete: :cascade }
      t.string :title
      t.string :given_name
      t.string :middle_name
      t.string :family_name
      t.string :suffix
      t.string :nickname
      t.string :gender
      t.string :pronouns
      t.string :birthdate
      t.boolean :birthdate_has_year
      t.boolean :birthdate_age_based

      t.timestamps
    end
  end
end
