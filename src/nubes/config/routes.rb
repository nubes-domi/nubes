Rails.application.routes.draw do
  root to: "home#index"

  ######
  # Signing in / Changing accounts
  ######

  # Identification (username / email)
  get "signin", to: "sessions#new"
  post "signin", to: "sessions#create"

  # Authentication (password / webauthn)
  get "signin/:method", to: "sessions#show", as: :authentication
  post "signin/:method", to: "sessions#update"

  ######
  # User management
  ######

  resource :me, controller: :profile, only: %i[show update] do
    get :name, as: :name
    get :birthdate, as: :birthdate
    get :gender, as: :gender
    get :qr, as: :qr
  end

  ######
  # Contacts
  ######

  resources :contacts do
    resources :contact_addresses, shallow: true, except: :index
    resources :contact_postal_addresses, shallow: true, except: :index
  end
end
