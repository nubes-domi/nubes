class NubesSchema < GraphQL::Schema
  mutation(MutationType)
  query(QueryType)

  # For batch-loading (see https://graphql-ruby.org/dataloader/overview.html)
  # use GraphQL::Dataloader

  def self.resolve_type(_type, obj, _ctx)
    obj.class.graphql_type
  end

  # def self.id_from_object(object, _type_definition, _query_ctx)
  #   object.id
  # end

  def self.object_from_id(id, _query_ctx)
    PrettyId.find(id)
  end
end
