class CreateUserSessions < ActiveRecord::Migration[7.0]
  def change
    create_table :user_sessions, id: :string do |t|
      t.references :user, null: false, type: :string, foreign_key: { on_delete: :cascade }
      t.datetime :expires_at, null: false
      t.string :ip_address
      t.string :user_agent

      t.timestamps
    end
  end
end
