require "rails_helper"

RSpec.describe "GraphQL" do
  fixtures :users, :contacts

  def run(variables = nil)
    NubesSchema.execute(query, variables:, context: { current_user: users(:joe) })
  end

  describe "Relay queries" do
    describe "node" do
      let(:query) do
        <<-GRAPHQL
          query($id: ID!) {
            node(id: $id) {
              id
            }
          }
        GRAPHQL
      end

      it "fetches a generic node" do
        contact = Contact.create!(given_name: "Mario", user_id: users(:joe).id)

        result = run("id" => contact.id)

        expect(result.dig("data", "node", "id")).to eq contact.id
      end
    end

    describe "nodes" do
      let(:query) do
        <<-GRAPHQL
          query($id: ID!) {
            nodes(ids: [$id]) {
              id
            }
          }
        GRAPHQL
      end

      it "fetches a list of generic nodes" do
        contact = Contact.create!(given_name: "Mario", user_id: users(:joe).id)

        result = run("id" => contact.id)

        expect(result.dig("data", "nodes", 0, "id")).to eq contact.id
      end
    end
  end
end
