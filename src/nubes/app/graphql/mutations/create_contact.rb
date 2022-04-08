module Mutations
  class CreateContact < Mutations::BaseMutation
    null true

    Contact.define_graphql_mutation(self, Types::Models::ContactType)

    def resolve(**params)
      op = Contacts::Operations::Create.call(current_user: context[:current_user], params:)
      if op.success?
        { contact: op["model"] }
      else
        { errors: to_errors(op["errors"]) }
      end
    end
  end
end
