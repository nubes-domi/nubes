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
end
