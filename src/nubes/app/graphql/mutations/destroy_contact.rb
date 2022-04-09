module Mutations
  class DestroyContact < Mutations::BaseMutation
    null true

    argument :id, ID, description: "ID of the contact to be destroyed"

    field :errors, [Types::Error]

    def resolve(id:)
      result = Contacts::Destroy.call(id:, user: context[:current_user])
      handle_failures(result) do
        {}
      end
    end
  end
end
