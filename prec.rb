require 'sinatra/base'
require 'sinatra/sequel'
require 'json'
require 'csv'
require 'logger'

require 'pry-remote'

class Prec < Sinatra::Base
  register Sinatra::SequelExtension

  configure do
    $stdout.sync = true
    database.logger = Logger.new($stdout)
    database.sql_log_level = :debug
  end

  get '/' do
    first_account_id = database[:accounts].select(:id).order(:id).limit(1).first[:id]
    redirect to("/accounts/#{first_account_id}/txns.html")
  end

  get '/accounts/:account_id/txns.html' do |account_id|
    @account_id = account_id.to_i
    @accounts = database[:accounts].order(:id)
    @txns = database[:txns].where(:account_id => account_id).
      reverse_order(:date, :id)
    @categories = database[:categories].order(:name).all
    if params[:filter] == 'uncategorised'
      @txns = @txns.where('category_id IS NULL')
    end
    erb :txns
  end

  post '/accounts/:account_id/txns.html' do |account_id|
    account_id = account_id.to_i
    tempfile = params[:txn_file][:tempfile]
    txns = []
    CSV.foreach(tempfile, :headers => :first_row) do |row|
      hash = hash_txn_row(row.to_a.map { |x| x.last })
      next unless
        database[:txns].where(:account_id => account_id, :hash => hash).empty?
      date = Time.parse(row[0]).iso8601
      description = row[1]
      debit_cents = parse_currency_string(row[2], true)
      credit_cents = parse_currency_string(row[3], true)
      balance_cents = parse_currency_string(row[4])
      database[:txns].insert(:account_id => account_id, :hash => hash,
        :date => date, :description => description,
        :debit_cents => debit_cents, :credit_cents => credit_cents,
        :balance_cents => balance_cents)
    end
    redirect to("/accounts/#{account_id}/txns.html")
  end

  patch '/txns/:txn_id.json' do |txn_id|
    txn_id = txn_id.to_i
    data = JSON.parse(request.body.read)
    database[:txns].where(:id => txn_id).
      update(:category_id => data['txn']['category_id'])
    200
  end

  get '/accounts.json' do
    @accounts = database[:accounts].order(:id)
    yajl :accounts
  end

  post '/accounts.json' do
    data = JSON.parse(request.body.read)
    database[:accounts].insert(:name => data['account']['name'])
    200
  end

  def hash_txn_row(row)
    concat_row = row.join('')
    hash = Digest::SHA2.new(512)
    hash.update(concat_row)
    hash.hexdigest
  end

  def parse_currency_string(string, force_positive = false)
    return nil unless string
    value = string.gsub!('.', '').to_i
    if force_positive and value < 0
      value * -1
    else
      value
    end
  end

  def format_currency(cents)
    return '0.00' if cents.nil?
    dollars = cents / 100.0
    "%.2f" % dollars
  end
end
