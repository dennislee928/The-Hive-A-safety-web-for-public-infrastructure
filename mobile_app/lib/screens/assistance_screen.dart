import 'package:flutter/material.dart';
import 'package:erh_safety_app/services/api_service.dart';
import 'package:erh_safety_app/services/storage_service.dart';

class AssistanceScreen extends StatefulWidget {
  const AssistanceScreen({super.key});

  @override
  State<AssistanceScreen> createState() => _AssistanceScreenState();
}

class _AssistanceScreenState extends State<AssistanceScreen> {
  final _formKey = GlobalKey<FormState>();
  final _descriptionController = TextEditingController();
  
  String? _selectedZone;
  String? _selectedRequestType;
  String? _selectedUrgency;
  bool _isSubmitting = false;
  
  final List<String> _zones = ['Z1', 'Z2', 'Z3', 'Z4'];
  final List<String> _requestTypes = ['medical', 'security', 'other'];
  final List<String> _urgencyLevels = ['high', 'medium', 'low'];
  
  @override
  void initState() {
    super.initState();
    _loadCurrentZone();
  }
  
  void _loadCurrentZone() {
    final currentZone = StorageService.getCurrentZone();
    setState(() {
      _selectedZone = currentZone ?? 'Z1';
    });
  }
  
  @override
  void dispose() {
    _descriptionController.dispose();
    super.dispose();
  }
  
  Future<void> _submitAssistanceRequest() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }
    
    if (_selectedZone == null || _selectedRequestType == null || _selectedUrgency == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please fill in all required fields')),
      );
      return;
    }
    
    setState(() {
      _isSubmitting = true;
    });
    
    final response = await ApiService.requestAssistance(
      zoneId: _selectedZone!,
      requestType: _selectedRequestType!,
      urgency: _selectedUrgency!,
      description: _descriptionController.text.trim(),
    );
    
    if (!mounted) return;
    
    setState(() {
      _isSubmitting = false;
    });
    
    if (response.success) {
      _descriptionController.clear();
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Assistance request submitted successfully'),
          backgroundColor: Colors.green,
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(response.error ?? 'Failed to submit assistance request'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Request Assistance'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Text(
                'Request immediate assistance',
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
                },
                validator: (value) {
                  if (value == null) {
                    return 'Please select a zone';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              DropdownButtonFormField<String>(
                value: _selectedRequestType,
                decoration: const InputDecoration(
                  labelText: 'Request Type',
                  border: OutlineInputBorder(),
                ),
                items: _requestTypes.map((type) {
                  return DropdownMenuItem(
                    value: type,
                    child: Text(type.toUpperCase()),
                  );
                }).toList(),
                onChanged: (value) {
                  setState(() {
                    _selectedRequestType = value;
                  });
                },
                validator: (value) {
                  if (value == null) {
                    return 'Please select a request type';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              DropdownButtonFormField<String>(
                value: _selectedUrgency,
                decoration: const InputDecoration(
                  labelText: 'Urgency',
                  border: OutlineInputBorder(),
                ),
                items: _urgencyLevels.map((urgency) {
                  Color? color;
                  if (urgency == 'high') {
                    color = Colors.red;
                  } else if (urgency == 'medium') {
                    color = Colors.orange;
                  } else {
                    color = Colors.green;
                  }
                  
                  return DropdownMenuItem(
                    value: urgency,
                    child: Row(
                      children: [
                        Container(
                          width: 12,
                          height: 12,
                          decoration: BoxDecoration(
                            color: color,
                            shape: BoxShape.circle,
                          ),
                        ),
                        const SizedBox(width: 8),
                        Text(urgency.toUpperCase()),
                      ],
                    ),
                  );
                }).toList(),
                onChanged: (value) {
                  setState(() {
                    _selectedUrgency = value;
                  });
                },
                validator: (value) {
                  if (value == null) {
                    return 'Please select urgency level';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _descriptionController,
                decoration: const InputDecoration(
                  labelText: 'Description (Optional)',
                  hintText: 'Provide additional details',
                  border: OutlineInputBorder(),
                ),
                maxLines: 4,
              ),
              const SizedBox(height: 32),
              ElevatedButton(
                onPressed: _isSubmitting ? null : _submitAssistanceRequest,
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  backgroundColor: _selectedUrgency == 'high' ? Colors.red : null,
                ),
                child: _isSubmitting
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Text('Request Assistance'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

