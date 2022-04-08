module Types
  class Error < Types::Base::Object
    field :code, String, null: false, description: "Machine readable error code"
    field :message, String, null: false, description: "Human readable message"
    field :path, [String], null: true, description: "Which argument caused the error"
  end
end
