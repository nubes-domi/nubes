module Types
  module Base
    class Connection < Types::Base::Object
      # add `nodes` and `pageInfo` fields, as well as `edge_type(...)` and `node_nullable(...)` overrides
      include GraphQL::Types::Relay::ConnectionBehaviors

      field :total_count, Integer, null: false

      def total_count
        object.items.size
      end
    end
  end
end
