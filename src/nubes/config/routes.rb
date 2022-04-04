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

  # Profile
  get "me", to: "profile#show"
  post "me", to: "profile#update"
  get "me/name", to: "profile#name", as: :me_name
  get "me/birthdate", to: "profile#birthdate", as: :me_birthdate
  get "me/gender", to: "profile#gender", as: :me_gender
  patch "me", to: "profile#update"
end
