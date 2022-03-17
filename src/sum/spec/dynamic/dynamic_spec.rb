require 'net/http'
require 'json'

describe 'Dynamic registration' do
  before(:all) do
    client_spec = {
      client_name: 'Nubes OIDC RSpec suite',
      grant_tyeps: ['authorization_code'],
      token_endpoint_auth_method: 'none',
      response_types: ['code'],
      redirect_uris: ['https://example.org/callback'],
      contacts: ['nubestest@example.com'],
      logo_uri: 'https://example.org/image.png',
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
    @registration_response_body = @registration_response.body
  end

  it 'responds with 201' do
    expect(@registration_response.code).to eq '201'
  end

  it 'responds with JSON' do
    expect(@registration_response['Content-Type']).to start_with 'application/json'
  end

  it 'does not error' do
    expect(@registration_response_body['error']).to be_nil
  end

  it 'provides client management credentials' do
    expect(@registration_response_body['registration_client_uri']).not_to be_nil
    expect(@registration_response_body['registration_access_token']).not_to be_nil
  end
end
