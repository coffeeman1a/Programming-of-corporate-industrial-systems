import unittest
from unittest.mock import mock_open, patch
from file_reader import read_file

class TestReadFile(unittest.TestCase):
    def test_valid_file_with_target(self):
        mock_data = 'I love coffee! Coffee is great! No tea, just coffee!'
        mock_target = 'coffee'
        with patch('builtins.open', mock_open(read_data=mock_data)), patch('os.path.exists', return_value=True):
            result = read_file('mock_file.txt', mock_target)
        self.assertEqual(result, (10, 3))
    
    def test_valid_file(self):
        mock_data = 'I love tea! Tea is great! Only tea!' 
        with patch('builtins.open', mock_open(read_data=mock_data)), patch('os.path.exists', return_value=True):
            result = read_file('mock_file.txt')
        self.assertEqual(result, (8, 0))

    def test_empty_file(self):
        mock_data = '' 
        with patch('builtins.open', mock_open(read_data=mock_data)), patch('os.path.exists', return_value=True):
            result = read_file('mock_file.txt')
        self.assertEqual(result, (0, 0))
    
    def test_file_not_found(self):
        with patch("builtins.open", side_effect=FileNotFoundError):
            with self.assertRaises(FileNotFoundError) as context: 
                read_file("non_existent_file.txt")
            self.assertIn("File not found:", str(context.exception))
    
    def test_empty_target(self):
        mock_data = 'I love coffee!' 
        with patch('builtins.open', mock_open(read_data=mock_data)), patch('os.path.exists', return_value=True):
            result = read_file('mock_file.txt', '')
        self.assertEqual(result, (3, 0))