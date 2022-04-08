module Types
  module Models
    class ContactAddressType < Types::Base::Object
      implements GraphQL::Types::Relay::Node

      description "Coordinates to communicate with your contacts (email, phone, irc)"

      ContactAddress.define_graphql_fields(self)
    end
  end
end
