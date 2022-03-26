Rails.application.routes.draw do
  ######
  # Signing in / Changing accounts
  ######

  # Sign in screen
  get "signin", to: "sessions#new"
  post "signin", to: "sessions#create"

  # Reauthentication
  get "reauth", to: "sessions#edit"
  post "reauth", to: "sessions#update"

  # Account switcher
  get "change-account", to: "switcher#index"
  post "change-account", to: "switcher#create"

  ######
  # OpenID Connect
  ######

  # Discovery
  get ".well-known/openid-configuration", to: "openid/discovery#show"
  namespace :openid do
    get "jwks", to: "discovery#keys"

    # Dynamic registration
    resources :clients, path: :registration, only: [:create, :show, :destroy]

    # Confirm grants
    get "authorization", to: "authorizations#create"
    get "authorization/:id", to: "authorizations#show"
    put "authorization/:id", to: "authorizations#update"

    post "token", to: "token#create"
    get "userinfo", to: "userinfo#show"
    post "userinfo", to: "userinfo#show"
  end
end
