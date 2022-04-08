module Mutations
  class DestroyContact < Mutations::BaseMutation
    null true

    argument :id, ID, description: "ID of the contact to be destroyed"

    field :errors, [Types::Error]

    def resolve(id:)
      op = Contacts::Operations::Destroy.call(id:, current_user: context[:current_user])
      if op.success?
        {}
      else
        { errors: to_errors(op["errors"]) }
      end
    end
  end
end
