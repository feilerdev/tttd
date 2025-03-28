defmodule HelloWorld do
  # TODO(tes.test): This implementation uses IO.puts which works only in terminal environments.
  def hello do
    IO.puts("Hello, World!")
  end
end

# Call the function
HelloWorld.hello()
