require 'securerandom'
require 'net/http'
require 'json'

def visit_authorization(
  client_id:,
  response_type: 'code',
  state: SecureRandom.urlsafe_base64,
  nonce: SecureRandom.urlsafe_base64,
  scope: 'openid',
  redirect_uri: 'https://example.org/callback'
)
  url = @discovery['authorization_endpoint']
  url += "?client_id=#{client_id}&response_type=#{response_type}&state=#{state}&nonce=#{nonce}&scope=#{scope}&redirect_uri=#{redirect_uri}"

  visit url
end

describe 'Authorization after dynamic registration', type: :feature do
  before(:all) do
    client_spec = {
      client_name: 'Nubes OIDC RSpec suite',
      grant_tyeps: ['authorization_code'],
      token_endpoint_auth_method: 'none',
      response_types: ['code'],
      redirect_uris: ['https://example.org/callback'],
      contacts: ['nubestest@example.com'],
      logo_uri: 'https://example.org/image.png',
      policy_uri: 'https://example.org/privacy',
      tos_uri: 'https://example.org/tos',
      jwks: {
        keys: []
      }
    }

    registration_endpoint = @discovery['registration_endpoint']

    uri = URI.parse(registration_endpoint)
    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = true
    request = Net::HTTP::Post.new(uri.request_uri, { 'Content-Type': 'text/json' })
    request.body = client_spec.to_json

    @registration_response = http.request(request)
    @registration_response_body = JSON.parse(@registration_response.body)
  end

  # it 'shows the provided custom logo' do
  #   visit_authorization(client_id: @registration_response_body['client_id'])
  #   expect(page.find('img')['src']).to have_content 'https://example.org/image.png'
  # end

  # it 'shows the provided privacy policy uri' do
  #   visit_authorization(client_id: @registration_response_body['client_id'])
  #   expect(page.find('a#policy_uri')['href']).to have_content 'https://example.org/privacy'
  # end

  # it 'shows the provided tos uri' do
  #   visit_authorization(client_id: @registration_response_body['client_id'])
  #   expect(page.find('a#tos_uri')['href']).to have_content 'https://example.org/tos'
  # end
end
