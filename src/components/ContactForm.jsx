import { useState } from 'react';

export default function ContactForm() {
  const [status, setStatus] = useState('idle');

  async function submit(event) {
    event.preventDefault();
    const form = event.currentTarget;
    const data = Object.fromEntries(new FormData(form));
    setStatus('sending');
    try {
      const response = await fetch('/api/contact', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!response.ok) throw new Error('Request failed');
      form.reset();
      setStatus('success');
    } catch {
      setStatus('error');
    }
  }

  const field = 'mt-2 w-full rounded-md border border-[#d6ddd8] bg-white px-4 py-3 text-sm text-ink outline-none focus:border-[#63892c]';
  return (
    <form onSubmit={submit} className="rounded-2xl bg-white p-6 shadow-xl shadow-[#1c3025]/5 sm:p-8">
      <div className="grid gap-5 sm:grid-cols-2">
        <label className="text-sm font-semibold">Name<input className={field} name="name" autoComplete="name" maxLength="100" required /></label>
        <label className="text-sm font-semibold">Work email<input className={field} name="email" type="email" autoComplete="email" maxLength="254" required /></label>
      </div>
      <label className="mt-5 block text-sm font-semibold">Organization<input className={field} name="organization" autoComplete="organization" maxLength="150" /></label>
      <label className="mt-5 block text-sm font-semibold">How can we help?<textarea className={`${field} min-h-32 resize-y`} name="message" maxLength="4000" required placeholder="Tell us a little about what’s on your mind." /></label>
      <button disabled={status === 'sending'} className="mt-5 flex w-full items-center justify-center gap-3 rounded-md bg-ink px-5 py-4 text-sm font-semibold text-white transition hover:bg-[#2b444b] disabled:opacity-60">{status === 'sending' ? 'Sending…' : 'Start the conversation'} <span aria-hidden="true">↗</span></button>
      <p className="mt-4 text-xs leading-5 text-[#6c7a77]">Please don’t include sensitive or confidential information in this form.</p>
      <p role="status" aria-live="polite" className="mt-3 text-sm">{status === 'success' ? 'Thank you. Your message has been sent.' : status === 'error' ? 'Your message could not be sent. Please try again later.' : ''}</p>
    </form>
  );
}
