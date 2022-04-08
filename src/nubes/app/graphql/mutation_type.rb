class MutationType < Types::Base::Object
  field :create_contact, mutation: Mutations::CreateContact
  field :update_contact, mutation: Mutations::UpdateContact
  field :destroy_contact, mutation: Mutations::DestroyContact

  field :create_contact_address, mutation: Mutations::CreateContactAddress
  field :update_contact_address, mutation: Mutations::UpdateContactAddress
  field :destroy_contact_address, mutation: Mutations::DestroyContactAddress

  field :create_contact_postal_address, mutation: Mutations::CreateContactPostalAddress
  field :update_contact_postal_address, mutation: Mutations::UpdateContactPostalAddress
  field :destroy_contact_postal_address, mutation: Mutations::DestroyContactPostalAddress
end
