require "rails_helper"

RSpec.describe "GraphQL" do
  fixtures :users, :contacts

  def run(variables = nil)
    NubesSchema.execute(query, variables:, context: { current_user: users(:joe) })
  end

  describe "Contacts" do
    describe "queries" do
      describe "contact" do
        let(:query) do
          <<-GRAPHQL
            query($id: ID!) {
              contact(id: $id) {
                id
              }
            }
          GRAPHQL
        end

        it "fetches a contact" do
          result = run("id" => contacts(:joe_contact1).id)
          expect(result.dig("data", "contact", "id")).not_to be_nil
        end
      end

      describe "contacts" do
        let(:query) do
          <<-GRAPHQL
            query {
              contacts {
                nodes {
                  id
                }
                totalCount
              }
            }
          GRAPHQL
        end

        it "fetches a list of contacts" do
          expect(run.dig("data", "contacts", "nodes", 0, "id")).not_to be_nil
        end

        it "can return total count of contacts" do
          expect(run.dig("data", "contacts", "totalCount")).not_to be_nil
        end
      end
    end
  end
end
