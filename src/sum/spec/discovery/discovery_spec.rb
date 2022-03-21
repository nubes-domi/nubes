require 'net/http'
require 'json'

describe 'Discovery endpoint' do
  it 'responds with 200' do
    expect(@discovery_response.code).to eq '200'
  end

  it 'responds with JSON' do
    expect(@discovery_response['Content-Type']).to start_with 'application/json'
  end

  it 'declares support for required response types' do
    expect(@discovery['response_types_supported']).to include 'code'
    expect(@discovery['response_types_supported']).to include 'id_token'
    expect(@discovery['response_types_supported']).to include 'token id_token'
  end

  it 'declares support for required grant types' do
    expect(@discovery['grant_types_supported']).to include 'authorization_code'
    expect(@discovery['grant_types_supported']).to include 'implicit'
  end

  it 'declares support for scopes' do
    expect(@discovery['scopes_supported']).to include 'openid'
  end

  it 'declares the issuer consistently with the expectations' do
    expect(@discovery['issuer']).to eq ENV['ISSUER']
  end

  it 'declares support for at least one of public and pairwise subject types' do
    expect(@discovery['subject_types_supported'] & %w[public pairwise]).to_not be_empty
  end

  it 'declares support for RS256 as id_token_signing_alg_values_supported' do
    expect(@discovery['id_token_signing_alg_values_supported']).to include('RS256')
  end

  it 'declares support for the "none" token endpoint auth method' do
    expect(@discovery['token_endpoint_auth_methods_supported']).to include 'none'
  end

  it 'declares support for request parameter' do
    expect(@discovery['request_parameter_supported']).to be true
  end

  it 'declares support for request_uri parameter' do
    expect(@discovery['request_uri_parameter_supported']).to be true
  end

  it 'declares support for RS256 signed request objects' do
    expect(@discovery['request_object_signing_alg_values_supported']).to include 'RS256'
  end

  xit 'declares support for the claims parameter' do
    expect(@discovery['claims_parameter_supported']).to be true
  end

  it 'declares supported claims' do
    expect(@discovery['claims_supported']).not_to be_empty
  end

  it 'only publishes https urls' do
    %w[authorization token userinfo registration].each do |key|
      expect(URI.parse(@discovery["#{key}_endpoint"]).scheme).to eq 'https'
    end
  end
end
