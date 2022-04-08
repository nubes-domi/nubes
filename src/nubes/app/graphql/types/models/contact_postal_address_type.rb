module Types
  module Models
    class ContactPostalAddressType < Types::Base::Object
      implements GraphQL::Types::Relay::Node

      description "A real world mailing address (street, city, etc) attached to a contact"

      ContactPostalAddress.define_graphql_fields(self)
    end
  end
end
