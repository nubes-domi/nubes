class ContactAddress < ApplicationRecord
  include PrettyId
  self.id_prefix = :caddr

  belongs_to :contact

  include Graphqlable

  field :name, String, description: "Name for this address. Eg: home, work, beach house."
  field :kind, String, description: "What type of contact address is this? phone, skype, email"
  field :data, String, description: "Value of the address: +44123456789, me@example.com, l33t.supa.h4x04"

  field :contact_id, GraphQL::Types::ID, only_create: true
  field :contact, Contact, readonly: true
end
