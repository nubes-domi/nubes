require 'net/http'
require 'json'

describe 'JWKs endpoint' do
  before(:all) do
    @jwks_response = Net::HTTP.get_response(URI.parse(@discovery['jwks_uri']))
    @jwks_document = JSON.parse(@jwks_response.body)
  end

  it 'responds with 200' do
    expect(@jwks_response.code).to eq '200'
  end

  it 'responds with JSON' do
    expect(@jwks_response['Content-Type']).to start_with 'application/json'
  end

  it 'does not publish symmetric keys' do
    expect(@jwks_document['keys'].select { |key| key['kty'] == 'oct' }).to be_empty
  end

  it 'does not publish private keys' do
    expect(@jwks_document['keys'].select do |key|
      (key.keys & %w[d p q dp dq qi oth r t k]).any?
    end).to be_empty
  end
end
