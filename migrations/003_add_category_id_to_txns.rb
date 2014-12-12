Sequel.migration do
  change do
    alter_table :txns do
      add_foreign_key :category_id, :categories
    end
  end
end
