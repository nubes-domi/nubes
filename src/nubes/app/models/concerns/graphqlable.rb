module Graphqlable
  extend ActiveSupport::Concern

  module ClassMethods
    attr_accessor :fields

    def field(name, type, **kwargs)
      self.fields ||= []
      self.fields << {
        name:, type:, options: kwargs
      }
    end

    def define_graphql_fields(klass)
      klass.field :id, GraphQL::Types::ID, null: false
      klass.field :created_at, GraphQL::Types::ISO8601DateTime,
                  null: false,
                  description: "When was this entity created."
      klass.field :updated_at, GraphQL::Types::ISO8601DateTime,
                  null: false,
                  description: "When was this entity last updated."

      self.fields.each do |f|
        klass.field f[:name], type_for(f[:type]), **f[:options].except(:readonly, :only_create)
      end
    end

    def define_graphql_mutation(klass, result_type, type: :create)
      klass.field result_type.graphql_name.underscore.to_sym, result_type
      klass.field :errors, [Types::Error]

      klass.argument :id, GraphQL::Types::ID, description: "ID of the entity to be updated" if type == :update

      fields_for_mutation(type).each do |f|
        klass.argument f[:name], type_for(f[:type]), required: false, **f[:options].except(:readonly, :only_create)
      end
    end

    def graphql_type
      Types::Models.const_get("#{self}Type")
    end

    private

    def fields_for_mutation(type)
      self.fields.reject do |field|
        field[:options][:readonly] || (type != :create && field[:options][:only_create])
      end
    end

    def type_for(type)
      if type.include?(Graphqlable)
        type.graphql_type
      elsif type.is_a?(Array) && type.length == 1 && type[0].include?(Graphqlable)
        type[0].graphql_type.connection_type
      elsif type.is_a?(Array)
        type.map { |field| type_for(field) }
      else
        type
      end
    end
  end
end
