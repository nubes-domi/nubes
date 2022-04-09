class MutationType < Types::Base::Object
  field :add_contact, mutation: Mutations::AddContact
  field :update_contact, mutation: Mutations::UpdateContact
  field :destroy_contact, mutation: Mutations::DestroyContact

  field :add_contact_address, mutation: Mutations::AddContactAddress
  field :update_contact_address, mutation: Mutations::UpdateContactAddress
  field :destroy_contact_address, mutation: Mutations::DestroyContactAddress

  field :add_contact_postal_address, mutation: Mutations::AddContactPostalAddress
  field :update_contact_postal_address, mutation: Mutations::UpdateContactPostalAddress
  field :destroy_contact_postal_address, mutation: Mutations::DestroyContactPostalAddress
end
