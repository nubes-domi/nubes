module Mutations
  class UpdateContact < Mutations::BaseMutation
    null true

    Contact.define_graphql_mutation(self, Types::Models::ContactType, type: :update)

    def resolve(id:, **attributes)
      result = Contacts::Update.call(user: context[:current_user], id:, attributes:)
      handle_failures(result) do |contact|
        { contact: }
      end
    end
  end
end
