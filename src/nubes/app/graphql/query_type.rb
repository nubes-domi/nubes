class QueryType < Types::Base::Object
  include GraphQL::Types::Relay::HasNodeField
  include GraphQL::Types::Relay::HasNodesField

  field :contacts, Types::Models::ContactType.connection_type, "List all your contacts"
  def contacts
    context[:current_user].contacts
  end

  field :contact, Types::Models::ContactType, "Find a contact by ID" do
    argument :id, ID
  end

  def contact(id:)
    context[:current_user].contacts.find(id)
  end
end
