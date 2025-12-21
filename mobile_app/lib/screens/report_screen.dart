import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:erh_safety_app/providers/auth_provider.dart';
import 'package:erh_safety_app/services/api_service.dart';
import 'package:erh_safety_app/services/storage_service.dart';

class ReportScreen extends StatefulWidget {
  const ReportScreen({super.key});

  @override
  State<ReportScreen> createState() => _ReportScreenState();
}

class _ReportScreenState extends State<ReportScreen> {
  final _formKey = GlobalKey<FormState>();
  final _contentController = TextEditingController();
  String? _selectedZone;
  bool _isSubmitting = false;
  
  final List<String> _zones = ['Z1', 'Z2', 'Z3', 'Z4'];
  
  @override
  void initState() {
    super.initState();
    _loadCurrentZone();
  }
  
  Future<void> _loadCurrentZone() async {
    final currentZone = StorageService.getCurrentZone();
    setState(() {
      _selectedZone = currentZone ?? 'Z1';
    });
  }
  
  @override
  void dispose() {
    _contentController.dispose();
    super.dispose();
  }
  
  Future<void> _submitReport() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }
    
    if (_selectedZone == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please select a zone')),
      );
      return;
    }
    
    setState(() {
      _isSubmitting = true;
    });
    
    final response = await ApiService.submitReport(
      zoneId: _selectedZone!,
      content: _contentController.text.trim(),
    );
    
    if (!mounted) return;
    
    setState(() {
      _isSubmitting = false;
    });
    
    if (response.success) {
      _contentController.clear();
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Report submitted successfully'),
          backgroundColor: Colors.green,
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(response.error ?? 'Failed to submit report'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Submit Report'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Text(
                'Report a safety issue or incident',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 24),
              DropdownButtonFormField<String>(
                value: _selectedZone,
                decoration: const InputDecoration(
                  labelText: 'Zone',
                  border: OutlineInputBorder(),
                ),
                items: _zones.map((zone) {
                  return DropdownMenuItem(
                    value: zone,
                    child: Text('Zone $zone'),
                  );
                }).toList(),
                onChanged: (value) {
                  setState(() {
                    _selectedZone = value;
                  });
                  if (value != null) {
                    StorageService.saveCurrentZone(value);
                  }
                },
                validator: (value) {
                  if (value == null) {
                    return 'Please select a zone';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 24),
              TextFormField(
                controller: _contentController,
                decoration: const InputDecoration(
                  labelText: 'Description',
                  hintText: 'Describe the safety issue or incident',
                  border: OutlineInputBorder(),
                ),
                maxLines: 5,
                validator: (value) {
                  if (value == null || value.trim().isEmpty) {
                    return 'Please enter a description';
                  }
                  if (value.trim().length < 10) {
                    return 'Description must be at least 10 characters';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 32),
              ElevatedButton(
                onPressed: _isSubmitting ? null : _submitReport,
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                ),
                child: _isSubmitting
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Text('Submit Report'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

