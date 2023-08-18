export function formatDateTime(date: Date | string): string {
  let d = new Date(date);
  return d.toLocaleDateString() + " " + d.toLocaleTimeString();
}

export function formatInitials(name: string): string {
  const words = name.toUpperCase().split(" ");
  if (words.length < 1) return "?";

  if (words.length < 2) return words[0][0] ?? "?";

  return (words[0][0] ?? "?") + (words[1][0] ?? "");
}
