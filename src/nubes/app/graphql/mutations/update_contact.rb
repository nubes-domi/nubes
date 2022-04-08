module Mutations
  class UpdateContact < Mutations::BaseMutation
    null true

    Contact.define_graphql_mutation(self, Types::Models::ContactType, type: :update)

    def resolve(id:, **params)
      op = Contacts::Operations::Update.call(current_user: context[:current_user], id:, params:)
      if op.success?
        { contact: op["model"] }
      else
        { errors: to_errors(op["errors"]) }
      end
    end
  end
end
