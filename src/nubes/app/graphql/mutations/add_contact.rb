module Mutations
  class AddContact < Mutations::BaseMutation
    null true

    Contact.define_graphql_mutation(self, Types::Models::ContactType)

    def resolve(**attributes)
      result = Contacts::Add.call(user: context[:current_user], attributes:)
      handle_failures(result) do |contact|
        { contact: }
      end
    end
  end
end
