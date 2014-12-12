Sequel.migration do
  up do
    create_table :accounts do
      primary_key :id
      column :name, 'varchar(255)', :null => false
    end

    account_id = self[:accounts].insert(:name => 'Everyday Account')

    drop_index :txns, :hash

    alter_table :txns do
      add_foreign_key :account_id, :accounts
    end

    self[:txns].update(:account_id => account_id)

    add_index :txns, [:account_id, :hash], :unique => true
  end
end
