class ContactPostalAddress < ApplicationRecord
  include PrettyId
  self.id_prefix = :cpaddr

  belongs_to :contact

  include Graphqlable

  field :name, String
  field :street, String
  field :locality, String
  field :region, String
  field :postal_code, String
  field :country, String
  field :preferred, GraphQL::Types::Boolean

  field :contact_id, GraphQL::Types::ID, only_create: true
  field :contact, Contact, readonly: true
end
