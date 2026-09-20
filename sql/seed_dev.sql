-- Dev seed only. Run AFTER migrate-up, never on prod.
-- Rerunnable: ON CONFLICT DO NOTHING.
-- All passwords are bcrypt for "password123".

INSERT INTO users (id, email, full_name, password, balance, phone_number, address, role)
OVERRIDING SYSTEM VALUE VALUES
(1, 'admin@test.com', 'Admin Satu', '$2a$10$kVmD5dWld28GWK.l9.yC..5LFECSqG9R.HarJo.aNFZMV.Tc8OWHO', 0, '081200000001', 'Jakarta', 'admin'),
(2, 'user1@test.com', 'User Satu', '$2a$10$kVmD5dWld28GWK.l9.yC..5LFECSqG9R.HarJo.aNFZMV.Tc8OWHO', 2000000, '081200000002', 'Bandung', 'user'),
(3, 'user2@test.com', 'User Dua', '$2a$10$kVmD5dWld28GWK.l9.yC..5LFECSqG9R.HarJo.aNFZMV.Tc8OWHO', 5000000, '081200000003', 'Surabaya', 'user'),
(4, 'user3@test.com', 'User Tiga', '$2a$10$kVmD5dWld28GWK.l9.yC..5LFECSqG9R.HarJo.aNFZMV.Tc8OWHO', 1000000, '081200000004', 'Medan', 'user')
ON CONFLICT (id) DO NOTHING;

INSERT INTO vehicles (id, name, description, quantity, daily_rate, category)
OVERRIDING SYSTEM VALUE VALUES
(1, 'Toyota Avanza', 'Mobil MPV keluarga sejuta umat, irit dan lega untuk 7 penumpang.', 5, 350000, 'MPV'),
(2, 'Honda Brio Satya', 'City car lincah dan hemat, cocok untuk macet Jakarta.', 3, 250000, 'City Car'),
(3, 'Toyota Fortuner VRZ', 'SUV diesel tangguh untuk luar kota dan medan berat.', 2, 800000, 'SUV'),
(4, 'Toyota Innova Zenix', 'MPV premium, kabin luas dan nyaman.', 2, 600000, 'Premium MPV'),
(5, 'Mitsubishi Xpander', 'Keluarga, suspensi empuk, desain modern.', 4, 400000, 'MPV'),
(6, 'Honda HR-V', 'Compact SUV stylish untuk eksekutif muda.', 3, 500000, 'Compact SUV')
ON CONFLICT (id) DO NOTHING;

INSERT INTO bookings (id, user_id, vehicle_id, booking_start, booking_end, total_price, status)
OVERRIDING SYSTEM VALUE VALUES
(1, 2, 1, '2026-08-01 08:00:00', '2026-08-03 08:00:00', 700000, 'completed'),
(2, 3, 3, '2026-08-15 10:00:00', '2026-08-20 10:00:00', 4000000, 'confirmed'),
(3, 4, 2, '2026-08-10 14:00:00', '2026-08-13 14:00:00', 750000, 'cancelled'),
(4, 2, 4, '2026-08-25 07:00:00', '2026-08-26 07:00:00', 600000, 'confirmed'),
(5, 4, 5, '2026-07-20 09:00:00', '2026-07-24 09:00:00', 1600000, 'completed'),
(6, 3, 6, '2026-08-18 08:00:00', '2026-08-20 08:00:00', 1000000, 'confirmed')
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT MAX(id) FROM users));
SELECT setval(pg_get_serial_sequence('vehicles', 'id'), (SELECT MAX(id) FROM vehicles));
SELECT setval(pg_get_serial_sequence('bookings', 'id'), (SELECT MAX(id) FROM bookings));
