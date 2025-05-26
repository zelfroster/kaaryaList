import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';

import Button from '../button.component'; // Adjust path as necessary

describe('Button Component', () => {
  test('renders with default props (primary variant)', () => {
    render(<Button value="Default Button" />);
    const buttonElement = screen.getByRole('button', { name: /Default Button/i });
    expect(buttonElement).toBeInTheDocument();
    expect(buttonElement).toHaveClass('bg-sky-500'); // Primary variant class
    expect(buttonElement).toHaveClass('px-4 py-2'); // Medium size, primary padding
  });

  test('renders with value prop', () => {
    const buttonText = 'Click Me';
    render(<Button value={buttonText} />);
    expect(screen.getByRole('button', { name: buttonText })).toBeInTheDocument();
  });

  test('renders with children prop for text content', () => {
    const childText = 'Child Text';
    render(<Button>{childText}</Button>);
    expect(screen.getByText(childText)).toBeInTheDocument();
  });

  test('renders an icon when passed via icon prop', () => {
    const Icon = () => <svg data-testid="test-icon" />;
    render(<Button icon={<Icon />} value="Button with Icon" />);
    expect(screen.getByTestId('test-icon')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Button with Icon/i })).toBeInTheDocument();
  });

  test('applies secondary variant classes', () => {
    render(<Button variant="secondary" value="Secondary Button" />);
    const buttonElement = screen.getByRole('button', { name: /Secondary Button/i });
    expect(buttonElement).toHaveClass('bg-zinc-700'); // Secondary variant class
    expect(buttonElement).toHaveClass('text-zinc-100');
    expect(buttonElement).toHaveClass('px-4 py-2'); // Medium size, secondary padding
  });

  test('applies icon variant classes and default padding', () => {
    const Icon = () => <svg data-testid="icon-variant-icon" />;
    render(<Button variant="icon" icon={<Icon />} />);
    const buttonElement = screen.getByRole('button');
    expect(buttonElement).toHaveClass('bg-transparent'); // Icon variant class
    expect(buttonElement).toHaveClass('text-zinc-400');
    expect(buttonElement).toHaveClass('p-2'); // Medium size, icon padding
    expect(buttonElement).toHaveClass('min-w-10 min-h-10');
  });

  test('handles click event', () => {
    const handleClick = jest.fn();
    render(<Button onClick={handleClick} value="Clickable" />);
    const buttonElement = screen.getByRole('button', { name: /Clickable/i });
    fireEvent.click(buttonElement);
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  test('is disabled when disabled prop is true', () => {
    const handleClick = jest.fn();
    render(<Button onClick={handleClick} value="Disabled Button" disabled />);
    const buttonElement = screen.getByRole('button', { name: /Disabled Button/i });
    
    expect(buttonElement).toBeDisabled();
    expect(buttonElement).toHaveClass('opacity-50 cursor-not-allowed');

    fireEvent.click(buttonElement);
    expect(handleClick).not.toHaveBeenCalled();
  });

  test('applies extraClassProps', () => {
    const extraClasses = 'my-custom-class another-class';
    render(<Button value="Extra Classes" extraClassProps={extraClasses} />);
    const buttonElement = screen.getByRole('button', { name: /Extra Classes/i });
    extraClasses.split(' ').forEach(cls => {
      expect(buttonElement).toHaveClass(cls);
    });
  });

  test('renders icon and value together', () => {
    const ButtonIcon = () => <span data-testid="button-icon">Icon</span>;
    render(<Button icon={<ButtonIcon />} value="Icon and Text" />);
    
    const buttonElement = screen.getByRole('button', { name: /Icon and Text/i });
    expect(buttonElement).toBeInTheDocument();
    expect(screen.getByTestId('button-icon')).toBeInTheDocument();
    expect(screen.getByText('Icon and Text')).toBeInTheDocument(); // Value is part of the accessible name
  });
  
  test('renders icon and children together', () => {
    const ButtonIcon = () => <span data-testid="button-icon-child">Icon</span>;
    render(<Button icon={<ButtonIcon />}>Child Content</Button>);
    
    const buttonElement = screen.getByRole('button'); // Name might be complex with icon + children
    expect(buttonElement).toBeInTheDocument();
    expect(screen.getByTestId('button-icon-child')).toBeInTheDocument();
    expect(screen.getByText('Child Content')).toBeInTheDocument();
  });

  test('default type is "button"', () => {
    render(<Button value="Test Type" />);
    const buttonElement = screen.getByRole('button', { name: /Test Type/i });
    expect(buttonElement).toHaveAttribute('type', 'button');
  });

  test('applies specified type attribute', () => {
    render(<Button value="Submit Type" type="submit" />);
    const buttonElement = screen.getByRole('button', { name: /Submit Type/i });
    expect(buttonElement).toHaveAttribute('type', 'submit');
  });
});
