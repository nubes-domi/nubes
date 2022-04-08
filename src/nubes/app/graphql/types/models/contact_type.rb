module Types
  module Models
    class ContactType < Types::Base::Object
      implements GraphQL::Types::Relay::Node

      description "A contact in your address book"

      Contact.define_graphql_fields(self)
    end
  end
end
