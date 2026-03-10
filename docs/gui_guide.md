# GoSvelte GUI Design Guide: Analytical Bento

This document defines the visual language, layout rules, and UI components for the GoSvelte frontend. The goal is to create a **premium, easy-to-digest, and data-precise** interface.

---

## 1. Core Aesthetic: The Analytical Bento
The interface is organized into a modular grid of "Bento Boxes" (rounded cards). Each card has a specific purpose and provides a clear, focused insight.

### Layout Principles
- **Grid System:** 12-column responsive grid.
- **The Gap Rule:** A consistent `20px` gap (gutter) between all Bento cards.
- **The 8px Multiplier:** All margins, paddings, and height increments must be multiples of 8 (8, 16, 24, 32, 48, 64).
- **Page Margins:** 
    - Mobile: `20px`
    - Tablet/Desktop: `40px`

---

## 2. Visual Tokens (Design System)

### Colors
| Token | HEX | Usage |
| :--- | :--- | :--- |
| `color-bg` | `#F8FAFC` | Page background (Smoke white) |
| `color-card` | `#FFFFFF` | Card background (Pure white) |
| `color-border` | `#E2E8F0` | Subtle 1px strokes |
| `brand-purple-dark` | `#4B0082` | Primary brand base (Royal Purple) |
| `brand-purple-light` | `#9F60CC` | Active states and highlights |
| `brand-yellow-primary` | `#EBCF41` | Secondary accents (Amber) |
| `brand-yellow-dark` | `#DDA81E` | Focus accents and warnings (Mustard) |
| **Semantic** | | |
| `color-income` | `#10B981` | Positive cash flow (Emerald) |
| `color-expense` | `#F43F5E` | Deductions (Rose) |
| `color-text-main` | `#1E293B` | Primary headings and body |
| `color-text-muted` | `#64748B` | Labels and placeholders |

### Typography
- **Primary Font:** `Inter` or `Geist Sans` (San-serif)
    - *Usage:* Navigation, Headings, UI Labels.
- **Data Font:** `JetBrains Mono` or `IBM Plex Mono` (Monospaced)
    - *Usage:* All currency values, percentages, and table numbers.
    - *Reason:* Ensures decimals align vertically for easier scanning.

| Style | Size | Weight |
| :--- | :--- | :--- |
| **Heading 1** | `32px` | 700 (Bold) |
| **Card Title** | `18px` | 600 (Semi-bold) |
| **Body** | `14px` | 400 (Regular) |
| **Label** | `12px` | 500 (Medium) |
| **Currency** | `16px` | 500 (Medium) |

### Elevation & Corners
- **Card Radius:** `24px` (Large, modern, friendly).
- **Button Radius:** `12px` (Distinguishable from cards).
- **Shadows:** Soft, diffuse shadows only.
    - `shadow-sm: 0 4px 6px -1px rgb(0 0 0 / 0.05)`
- **Borders:** Every card has a `1px solid var(--color-border)`.

---

## 3. Interaction Standards (Micro-UX)

### 1. Feedback Loops
- **Active State:** Buttons should shrink slightly (`scale: 0.98`) when pressed.
- **Hover State:** Bento cards should lift by `2px` and gain a slightly darker border.
- **Loading:** Use **Skeleton Shimmer** effects that match the exact shape of the Bento cards. No generic spinning wheels.

### 2. Form Entry (The "Speed" Standard)
- **Mobile Centric:** Large touch targets (minimum `44px` height).
- **Auto-Focus:** When opening "Quick Add," automatically focus the amount input.
- **Numeric Keyboard:** Always trigger the decimal keyboard (`inputmode="decimal"`) for currency entries.

### 3. Transitions
- Use Svelte's `fly` and `fade` transitions for all modal and page entries.
- Duration: `200ms` - `300ms` (Fast but fluid).

---

## 4. Specific View Logic

### The Dashboard (Desktop)
- **Top Row:** Net Worth (Large Card), Debt Summary (Medium Card).
- **Middle Grid:** Cash Flow Chart (Large), Spending Breakdown (Medium).
- **Sidebar:** Recent Activity List.

### The Mobile View (Entry Hub)
- **Sticky Bottom Navigation:** Home, Transactions, Add (+), Budgets, Wishlist.
- **Add Button (+):** A floating purple button in the center of the dock.
- **Feed Layout:** Single column of cards.

---

## 5. CSS Implementation (Root Variables)

```css
:root {
  /* Spacing */
  --gap: 20px;
  --radius-card: 24px;
  --radius-btn: 12px;

  /* Typography */
  --font-sans: 'Inter', system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', monospace;

  /* Brand Colors */
  --purple-dark: #4b0082;
  --purple-light: #9f60cc;
  --yellow-primary: #ebcf41;
  --yellow-dark: #dda81e;

  /* Layout */
  --bg: #f8fafc;
  --card: #ffffff;
  --border: #e2e8f0;
  --text: #1e293b;
  --text-muted: #64748b;
  
  --income: #10b981;
  --expense: #f43f5e;
  --accent: var(--purple-dark);
}
```
