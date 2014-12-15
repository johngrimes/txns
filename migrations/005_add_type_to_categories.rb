Sequel.migration do
  up do
    alter_table :categories do
      add_column :category_type, Integer, :default => 1, :null => false
    end
  end
end
